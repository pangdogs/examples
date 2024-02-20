package main

import (
	"context"
	"encoding/json"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/plugins/broker/nats_broker"
	"git.golaxy.org/framework/plugins/discovery/cache_discovery"
	"git.golaxy.org/framework/plugins/discovery/redis_discovery"
	"git.golaxy.org/framework/plugins/dserv"
	"git.golaxy.org/framework/plugins/dsync/redis_dsync"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/log/console_log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Register("demo", DemoComp{})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle, console_log.With.Level(log.DebugLevel))
	nats_broker.Install(pluginBundle, nats_broker.With.CustomAddresses("192.168.10.5:4222"))
	cache_discovery.Install(pluginBundle, cache_discovery.With.Wrap(redis_discovery.NewRegistry(redis_discovery.With.CustomAddress("192.168.10.5:6379"))))
	redis_dsync.Install(pluginBundle, redis_dsync.With.CustomAddress("192.168.10.5:6379"), redis_dsync.With.CustomDB(1))
	dserv.Install(pluginBundle)

	// 创建服务上下文与服务，并开始运行
	<-core.NewService(service.NewContext(
		service.With.EntityLib(entityLib),
		service.With.PluginBundle(pluginBundle),
		service.With.Name("demo_dserv"),
		service.With.RunningHandler(generic.CastDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

			// 监听退出信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			// 监听消息
			dserv.Using(ctx).WatchMsg(context.Background(), generic.CastDelegateFunc2(
				func(topic string, mp gap.MsgPacket) error {
					data, _ := json.Marshal(mp)
					log.Infof(ctx, "receive => topic:%q, msg-packet:%s", topic, data)
					return nil
				},
			))

			go func() {
				<-sigChan
				ctx.GetCancelFunc()()
			}()

			// 创建运行时上下文与运行时，并开始运行
			rt := core.NewRuntime(
				runtime.NewContext(ctx,
					runtime.With.Context.RunningHandler(generic.CastDelegateAction2(func(_ runtime.Context, state runtime.RunningState) {
						if state != runtime.RunningState_Terminated {
							return
						}
						ctx.GetCancelFunc()()
					})),
				),
				core.With.Runtime.Frame(runtime.NewFrame()),
				core.With.Runtime.AutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			core.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
				entity, err := core.CreateEntity(ctx).
					Prototype("demo").
					Scope(ec.Scope_Global).
					Spawn()
				if err != nil {
					log.Panic(service.Current(ctx), err)
				}
				log.Infof(service.Current(ctx), "create entity %q finish", entity)
			})
		})),
	)).Run()
}
