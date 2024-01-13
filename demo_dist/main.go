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
	"git.golaxy.org/plugins/broker/nats_broker"
	"git.golaxy.org/plugins/dist"
	"git.golaxy.org/plugins/dsync/redis_dsync"
	"git.golaxy.org/plugins/gap"
	"git.golaxy.org/plugins/log"
	"git.golaxy.org/plugins/log/console_log"
	"git.golaxy.org/plugins/registry/cache_registry"
	"git.golaxy.org/plugins/registry/redis_registry"
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
	console_log.Install(pluginBundle, console_log.Option{}.Level(log.InfoLevel))
	nats_broker.Install(pluginBundle, nats_broker.Option{}.FastAddresses("127.0.0.1:4222"))
	cache_registry.Install(pluginBundle, cache_registry.Option{}.Wrap(redis_registry.NewRegistry(redis_registry.Option{}.FastAddress("127.0.0.1:6379"))))
	redis_dsync.Install(pluginBundle, redis_dsync.Option{}.FastAddress("127.0.0.1:6379"), redis_dsync.Option{}.FastDB(1))
	dist.Install(pluginBundle)

	// 创建服务上下文与服务，并开始运行
	<-core.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_dist"),
		service.Option{}.RunningHandler(generic.CastDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

			// 监听退出信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			// 监听消息
			dist.WatchMsg(ctx, context.Background(), generic.CastDelegateFunc2(
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
					runtime.Option{}.Context.RunningHandler(generic.CastDelegateAction2(func(_ runtime.Context, state runtime.RunningState) {
						if state != runtime.RunningState_Terminated {
							return
						}
						ctx.GetCancelFunc()()
					})),
				),
				core.Option{}.Runtime.Frame(runtime.NewFrame()),
				core.Option{}.Runtime.AutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			core.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
				entity, err := core.CreateEntity(ctx,
					core.Option{}.EntityCreator.Prototype("demo"),
					core.Option{}.EntityCreator.Scope(ec.Scope_Global),
				).Spawn()
				if err != nil {
					log.Panic(service.Current(ctx), err)
				}
				log.Infof(service.Current(ctx), "create entity %q finish", entity)
			})
		})),
	)).Run()
}
