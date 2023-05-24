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
	redis_dsync "kit.golaxy.org/plugins/dsync/redis"
	"kit.golaxy.org/plugins/logger"
	zap_logger "kit.golaxy.org/plugins/logger/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("demo", []string{
		defineDemoComp.Implementation,
	})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()

	// 安装日志插件
	zapLogger, _ := zap_logger.NewConsoleZapLogger(zapcore.DebugLevel, "\t", "", 0, true, true)
	zap_logger.Install(pluginBundle, zap_logger.WithOption{}.ZapLogger(zapLogger), zap_logger.WithOption{}.Fields(0))

	// 安装redis服务发现插件
	redis_dsync.Install(pluginBundle, redis_dsync.WithOption{}.FastAddress("127.0.0.1:6379"), redis_dsync.WithOption{}.FastDBIndex(0))

	//// 安装etcd分布式同步插件
	//etcd_dsync.Install(pluginBundle, etcd_dsync.WithOption{}.FastAddresses("127.0.0.1:12379"))

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.WithOption{}.EntityLib(entityLib),
		service.WithOption{}.PluginBundle(pluginBundle),
		service.WithOption{}.Name("demo_registry"),
		service.WithOption{}.StartedCallback(func(serviceCtx service.Context) {
			// 监听退出信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			go func() {
				<-sigChan
				serviceCtx.GetCancelFunc()()
			}()

			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(runtime.NewContext(serviceCtx),
				golaxy.WithOption{}.RuntimeFrame(runtime.NewFrame(30, 0, false)),
				golaxy.WithOption{}.RuntimeAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			golaxy.AsyncVoid(rt, func(runtimeCtx runtime.Context) {
				_, err := golaxy.NewEntityCreator(runtimeCtx,
					pt.WithOption{}.Prototype("demo"),
					pt.WithOption{}.Scope(ec.Scope_Global),
				).Spawn()
				if err != nil {
					logger.Panic(service.Get(runtimeCtx), err)
				}
			})
		}),
	)).Run()
}
