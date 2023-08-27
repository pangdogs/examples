package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/uid"
	"kit.golaxy.org/plugins/gtp_gate"
	"kit.golaxy.org/plugins/logger"
)

// SessionStateChangedHandler 会话状态变化的处理器
func SessionStateChangedHandler(session gtp_gate.Session, old, new gtp_gate.SessionState) {
	id, err := uid.UnmarshalText([]byte(session.GetId()))
	if err != nil {
		logger.Panic(session.GetContext(), err)
	}

	switch new {
	case gtp_gate.SessionState_Confirmed:
		// 创建运行时上下文与运行时，并开始运行
		rt := golaxy.NewRuntime(runtime.NewContext(session.GetContext(), runtime.Option{}.Context.Context(session)),
			golaxy.Option{}.Runtime.AutoRun(true),
		)

		// 在运行时线程环境中，创建实体
		<-golaxy.AsyncVoid(rt, func(runtimeCtx runtime.Context) {
			_, err := golaxy.EntityCreator{Context: runtimeCtx}.Clone().
				Options(
					golaxy.Option{}.EntityCreator.Prototype("demo"),
					golaxy.Option{}.EntityCreator.Scope(ec.Scope_Global),
					golaxy.Option{}.EntityCreator.EntityConstructor(func(entity ec.Entity) {
						pt.Cast[IDemoComp](entity).(IDemoCompConstructor).Constructor(session)
					}),
					golaxy.Option{}.EntityCreator.PersistId(id),
				).Spawn()
			if err != nil {
				logger.Panic(service.Current(runtimeCtx), err)
			}
		})

	case gtp_gate.SessionState_Death:
		session.GetContext().AsyncCallVoid(id, func(entity ec.Entity) {
			entity.DestroySelf()
		})
	}
}
