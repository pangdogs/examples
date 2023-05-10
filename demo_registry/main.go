package main

import (
	"go.uber.org/zap/zapcore"
	_ "kit.golaxy.org/components"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/logger"
	zap_logger "kit.golaxy.org/plugins/logger/zap"
	cache_registry "kit.golaxy.org/plugins/registry/cache"
	_ "kit.golaxy.org/plugins/registry/etcd"
	redis_registry "kit.golaxy.org/plugins/registry/redis"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("demo", []string{
		defineDemoComp.Path,
	})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()

	zapLogger, _ := zap_logger.NewZapConsoleLogger(zapcore.DebugLevel, "\t", "", 0, true, true)
	zap_logger.Install(pluginBundle, zap_logger.WithZapOption{}.ZapLogger(zapLogger), zap_logger.WithZapOption{}.Fields(0))

	//r := etcd_registry.NewEtcdRegistry(etcd_registry.WithEtcdOption{}.FastAddresses("localhost:2379"))
	r := redis_registry.NewRedisRegistry(redis_registry.WithRedisOption{}.FastAddress("localhost:6379"))
	cache_registry.Install(pluginBundle, cache_registry.WithCacheOption{}.Cached(r))

	// 创建服务上下文
	ctx := service.NewContext(
		service.WithContextOption{}.EntityLib(entityLib),
		service.WithContextOption{}.PluginBundle(pluginBundle),
		service.WithContextOption{}.Name("demo_registry"),
		service.WithContextOption{}.StartedCallback(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.WithContextOption{}.AutoRecover(false),
				),
				golaxy.WithRuntimeOption{}.Frame(runtime.NewFrame(30, 0, false)),
				golaxy.WithRuntimeOption{}.EnableAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			rt.GetContext().AsyncCallNoRet(func() {
				entity, err := golaxy.NewEntityCreator(rt.GetContext(),
					pt.WithEntityOption{}.Prototype("demo"),
					pt.WithEntityOption{}.Scope(ec.Scope_Global),
				).Spawn()
				if err != nil {
					logger.Panic(rt.GetContext(), err)
				}

				logger.Debugf(rt.GetContext(), "create entity %q finish", entity)
			})
		}),
	)

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		ctx.GetCancelFunc()()
	}()

	// 开始运行
	<-golaxy.NewService(ctx).Run()
}
