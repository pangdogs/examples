package main

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/plugins/discovery/cache_discovery"
	"git.golaxy.org/plugins/discovery/etcd_discovery"
	"git.golaxy.org/plugins/discovery/redis_discovery"
	"git.golaxy.org/plugins/log"
	"git.golaxy.org/plugins/log/console_log"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Register("demo", DemoComp{})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle)

	// 创建etcd服务发现插件
	etcdRegistry := etcd_discovery.NewRegistry(etcd_discovery.Option{}.FastAddresses("192.168.10.8:2379"))
	_ = etcdRegistry

	// 创建redis服务发现插件
	redisRegistry := redis_discovery.NewRegistry(redis_discovery.Option{}.FastAddress("192.168.10.8:6379"))
	_ = redisRegistry

	// 安装服务发现插件，使用服务缓存插件包装其他服务发现插件
	cache_discovery.Install(pluginBundle, cache_discovery.Option{}.Wrap(redisRegistry))

	// 创建服务上下文与服务，并开始运行
	<-core.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_discovery"),
		service.Option{}.RunningHandler(generic.CastDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

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
