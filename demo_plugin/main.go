package main

import (
	"go.uber.org/zap/zapcore"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/logger"
	zap_logger "kit.golaxy.org/plugins/logger/zap"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("demo", []string{
		defineDemoComp.Path,
	})

	// 创建插件库，安装插件
	pluginBundle := plugin.NewPluginBundle()
	zapLogger, _ := zap_logger.NewZapConsoleLogger(zapcore.DebugLevel, "\t", "", 0, true, true)
	zap_logger.Install(pluginBundle, zap_logger.WithZapOption{}.ZapLogger(zapLogger), zap_logger.WithZapOption{}.Fields(0))
	defineDemoPlugin.Install(pluginBundle)

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.WithContextOption{}.EntityLib(entityLib),
		service.WithContextOption{}.PluginBundle(pluginBundle),
		service.WithContextOption{}.Name("demo_plugin"),
		service.WithContextOption{}.StartedCallback(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.WithContextOption{}.StoppedCallback(func(runtime.Context) { serviceCtx.GetCancelFunc()() }),
					runtime.WithContextOption{}.AutoRecover(false),
				),
				golaxy.WithRuntimeOption{}.EnableAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			rt.GetContext().AsyncCallNoRet(func() {
				_, err := golaxy.NewEntityCreator(rt.GetContext(),
					pt.WithEntityOption{}.Prototype("demo"),
					pt.WithEntityOption{}.Scope(ec.Scope_Global),
				).Spawn()
				if err != nil {
					logger.Panic(service.Get(rt), err)
				}
			})
		}),
	)).Run()
}
