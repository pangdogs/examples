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
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("demo", []string{
		defineDemoComp.Implementation,
	})

	// 创建插件包
	pluginBundle := plugin.NewPluginBundle()

	// 安装日志插件
	zapLogger, _ := zap_logger.NewConsoleZapLogger(zapcore.DebugLevel, "\t", "", 0, true, true)
	zap_logger.Install(pluginBundle, zap_logger.WithOption{}.ZapLogger(zapLogger), zap_logger.WithOption{}.Fields(0))

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.WithOption{}.EntityLib(entityLib),
		service.WithOption{}.PluginBundle(pluginBundle),
		service.WithOption{}.Name("demo_ec"),
		service.WithOption{}.StartedCb(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.WithOption{}.StoppedCb(func(runtime.Context) { serviceCtx.GetCancelFunc()() }),
				),
				golaxy.WithOption{}.RuntimeFrame(runtime.NewFrame(30, 300, false)),
				golaxy.WithOption{}.RuntimeAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			golaxy.AsyncVoid(rt, func(runtimeCtx runtime.Context) {
				entity, err := golaxy.NewEntityCreator(runtimeCtx,
					pt.WithOption{}.Prototype("demo"),
					pt.WithOption{}.Scope(ec.Scope_Global),
				).Spawn()
				if err != nil {
					logger.Panic(service.Get(runtimeCtx), err)
				}

				logger.Debugf(service.Get(runtimeCtx), "create entity %q finish", entity)
			})
		}),
	)).Run()
}
