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
	"kit.golaxy.org/plugins/logger/zap_logger"
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
	zap_logger.Install(pluginBundle, zap_logger.Option{}.ZapLogger(zapLogger), zap_logger.Option{}.Fields(0))

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_ec"),
		service.Option{}.RunningHandler(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(ctx,
					runtime.Option{}.Context.RunningHandler(func(_ runtime.Context, state runtime.RunningState) {
						if state != runtime.RunningState_Terminated {
							return
						}
						ctx.GetCancelFunc()()
					}),
				),
				golaxy.Option{}.Runtime.Frame(runtime.NewFrame(runtime.Option{}.Frame.TotalFrames(300))),
				golaxy.Option{}.Runtime.AutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			golaxy.AsyncVoid(rt, func(runtimeCtx runtime.Context) {
				entity, err := golaxy.EntityCreator{Context: runtimeCtx}.Clone().
					Options(
						golaxy.Option{}.EntityCreator.Prototype("demo"),
						golaxy.Option{}.EntityCreator.Scope(ec.Scope_Global),
					).Spawn()
				if err != nil {
					logger.Panic(service.Current(runtimeCtx), err)
				}
				logger.Debugf(service.Current(runtimeCtx), "create entity %q finish", entity)
			})
		}),
	)).Run()
}
