package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util/generic"
	"kit.golaxy.org/plugins/log"
	"kit.golaxy.org/plugins/log/console_log"
	"kit.golaxy.org/plugins/registry/cache_registry"
	"kit.golaxy.org/plugins/registry/etcd_registry"
	"kit.golaxy.org/plugins/registry/redis_registry"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Register("demo", DemoComp{})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle)

	// 创建etcd服务发现插件
	etcdRegistry := etcd_registry.NewRegistry(etcd_registry.Option{}.FastAddresses("192.168.10.8:2379"))
	_ = etcdRegistry

	// 创建redis服务发现插件
	redisRegistry := redis_registry.NewRegistry(redis_registry.Option{}.FastAddress("192.168.10.8:6379"))
	_ = redisRegistry

	// 安装服务发现插件，使用服务缓存插件包装其他服务发现插件
	cache_registry.Install(pluginBundle, cache_registry.Option{}.Wrap(redisRegistry))

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_registry"),
		service.Option{}.RunningHandler(generic.CastDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

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
			})
		})),
	)).Run()
}
