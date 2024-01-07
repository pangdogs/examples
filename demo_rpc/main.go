package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util/generic"
	"kit.golaxy.org/plugins/broker/nats_broker"
	"kit.golaxy.org/plugins/distributed"
	"kit.golaxy.org/plugins/dsync/redis_dsync"
	"kit.golaxy.org/plugins/log"
	"kit.golaxy.org/plugins/log/console_log"
	"kit.golaxy.org/plugins/registry/cache_registry"
	"kit.golaxy.org/plugins/registry/redis_registry"
	"kit.golaxy.org/plugins/rpc"
	"kit.golaxy.org/plugins/util/concurrent"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var entities = concurrent.MakeLockedSlice[string](0, 0)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Register("demo",
		pt.CompWithName{
			Comp: DemoComp{},
			Name: "DemoComp",
		},
	)

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle, console_log.Option{}.Level(log.InfoLevel), console_log.Option{}.TimestampLayout(time.StampMilli))
	nats_broker.Install(pluginBundle, nats_broker.Option{}.FastAddresses("192.168.10.7:4222"))
	cache_registry.Install(pluginBundle, cache_registry.Option{}.Wrap(redis_registry.NewRegistry(redis_registry.Option{}.FastAddress("192.168.10.7:6379"))))
	redis_dsync.Install(pluginBundle, redis_dsync.Option{}.FastAddress("192.168.10.7:6379"), redis_dsync.Option{}.FastDB(1))
	distributed.Install(pluginBundle, distributed.Option{}.FutureTimeout(time.Minute))
	rpc.Install(pluginBundle)

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_rpc"),
		service.Option{}.RunningHandler(generic.CastDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

			// 监听退出信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			go func() {
				<-sigChan
				ctx.GetCancelFunc()()
			}()

			for i := 0; i < 3; i++ {
				// 创建运行时上下文与运行时，并开始运行
				rt := golaxy.NewRuntime(
					runtime.NewContext(ctx,
						runtime.Option{}.Context.RunningHandler(generic.CastDelegateAction2(func(_ runtime.Context, state runtime.RunningState) {
							if state != runtime.RunningState_Terminated {
								return
							}
							ctx.GetCancelFunc()()
						})),
					),
					golaxy.Option{}.Runtime.Frame(runtime.NewFrame()),
					golaxy.Option{}.Runtime.AutoRun(true),
				)

				// 在运行时线程环境中，创建实体
				golaxy.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
					entity, err := golaxy.CreateEntity(ctx,
						golaxy.Option{}.EntityCreator.Prototype("demo"),
						golaxy.Option{}.EntityCreator.Scope(ec.Scope_Global),
					).Spawn()
					if err != nil {
						log.Panic(service.Current(ctx), err)
					}
					log.Infof(service.Current(ctx), "create entity %q finish", entity)

					entities.Append(entity.GetId().String())
				})
			}
		})),
	)).Run()
}
