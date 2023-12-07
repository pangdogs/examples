package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util/generic"
	"kit.golaxy.org/golaxy/util/uid"
	"kit.golaxy.org/plugins/gtp_gate"
	"kit.golaxy.org/plugins/log"
	"kit.golaxy.org/plugins/log/console_log"
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
	console_log.Install(pluginBundle)

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
	<-golaxy.NewService(service.NewContext(
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

func handleSessionStateChanged(session gtp_gate.Session, old, new gtp_gate.SessionState) {
	switch new {
	case gtp_gate.SessionState_Confirmed:
		// 创建运行时上下文与运行时，并开始运行
		rt := golaxy.NewRuntime(runtime.NewContext(session.GetContext()),
			golaxy.Option{}.Runtime.AutoRun(true),
		)

		// 在运行时线程环境中，创建实体
		golaxy.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
			entity, err := golaxy.CreateEntity(ctx,
				golaxy.Option{}.EntityCreator.Prototype("demo"),
				golaxy.Option{}.EntityCreator.Scope(ec.Scope_Global),
				golaxy.Option{}.EntityCreator.PersistId(uid.Id(session.GetId())),
				golaxy.Option{}.EntityCreator.Meta(map[string]any{"session": session}),
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
