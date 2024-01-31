package main

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/core/util/uid"
	"git.golaxy.org/framework/plugins/gtp_gate"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/log/console_log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing endpoints")
	}

	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Register("demo", DemoComp{})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle, console_log.Option{}.Level(log.DebugLevel))

	// 安装网关插件
	gtp_gate.Install(pluginBundle,
		gtp_gate.Option{}.Gate.Endpoints(os.Args[1:]...),
		gtp_gate.Option{}.Gate.IOTimeout(3*time.Second),
		gtp_gate.Option{}.Gate.IOBufferCap(1024*1024*5),
		gtp_gate.Option{}.Gate.AgreeClientEncryptionProposal(true),
		gtp_gate.Option{}.Gate.AgreeClientCompressionProposal(true),
		gtp_gate.Option{}.Gate.CompressedSize(128),
		gtp_gate.Option{}.Gate.SessionInactiveTimeout(time.Hour),
		gtp_gate.Option{}.Gate.SessionStateChangedHandler(generic.CastDelegateAction3(handleSessionStateChanged)),
	)

	// 创建服务上下文与服务，并开始运行
	<-core.NewService(service.NewContext(
		service.Option{}.EntityLib(entityLib),
		service.Option{}.PluginBundle(pluginBundle),
		service.Option{}.Name("demo_gate"),
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
		})),
	)).Run()
}

func handleSessionStateChanged(session gtp_gate.ISession, old, new gtp_gate.SessionState) {
	switch new {
	case gtp_gate.SessionState_Confirmed:
		// 创建运行时上下文与运行时，并开始运行
		rt := core.NewRuntime(runtime.NewContext(session.GetContext()),
			core.Option{}.Runtime.AutoRun(true),
		)

		// 在运行时线程环境中，创建实体
		core.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
			entity, err := core.CreateEntity(ctx,
				core.Option{}.EntityCreator.Prototype("demo"),
				core.Option{}.EntityCreator.Scope(ec.Scope_Global),
				core.Option{}.EntityCreator.PersistId(uid.Id(session.GetId())),
				core.Option{}.EntityCreator.Meta(map[string]any{"session": session}),
			).Spawn()
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}
			log.Infof(service.Current(ctx), "create entity %q finish", entity)
		}).Wait(session.GetContext())

	case gtp_gate.SessionState_Death:
		session.GetContext().CallVoid(uid.Id(session.GetId()), func(entity ec.Entity, _ ...any) {
			entity.DestroySelf()
		})
	}
}
