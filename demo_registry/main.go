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
	etcd_registry "kit.golaxy.org/plugins/registry/etcd"
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

	zapLogger, _ := zap_logger.NewConsoleZapLogger(zapcore.InfoLevel, "\t", "", 0, true, true)
	zap_logger.Install(pluginBundle, zap_logger.WithOption{}.ZapLogger(zapLogger), zap_logger.WithOption{}.Fields(0))

	etcdRegistry := etcd_registry.NewRegistry(etcd_registry.WithOption{}.FastAddresses("DevelopWin10:2379"))
	_ = etcdRegistry
	redisRegistry := redis_registry.NewRegistry(redis_registry.WithOption{}.FastAddress("DevelopWin10:6379"), redis_registry.WithOption{}.FastDBIndex(10))
	_ = redisRegistry
	cache_registry.Install(pluginBundle, cache_registry.WithOption{}.Cached(redisRegistry))

	// 创建服务上下文
	ctx := service.NewContext(
		service.WithOption{}.EntityLib(entityLib),
		service.WithOption{}.PluginBundle(pluginBundle),
		service.WithOption{}.Name("demo_registry"),
		service.WithOption{}.StartedCallback(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.WithOption{}.AutoRecover(false),
				),
				golaxy.WithRuntimeOption{}.Frame(runtime.NewFrame(30, 0, false)),
				golaxy.WithRuntimeOption{}.EnableAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			rt.GetContext().AsyncCallNoRet(func() {
				entity, err := golaxy.NewEntityCreator(rt.GetContext(),
					pt.WithOption{}.Prototype("demo"),
					pt.WithOption{}.Scope(ec.Scope_Global),
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
