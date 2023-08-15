package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/pt"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/uid"
	"kit.golaxy.org/plugins/gate"
	"kit.golaxy.org/plugins/logger"
)

// SessionStateChangedHandler 会话状态变化的处理器
func SessionStateChangedHandler(session gate.Session, old, new gate.SessionState) {
	logger.Infof(session.GetContext(), "session %q state %q => %q", session.GetId(), old, new)

	var id uid.Id
	if err := id.UnmarshalText([]byte(session.GetId())); err != nil {
		logger.Panic(session.GetContext(), err)
	}
	
	switch new {
	case gate.SessionState_Confirmed:
		// 创建运行时上下文与运行时，并开始运行
		rt := golaxy.NewRuntime(runtime.NewContext(session.GetContext(), runtime.Option{}.Context(session)),
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
				logger.Panic(service.Get(runtimeCtx), err)
			}
		})

	case gate.SessionState_Death:
		session.GetContext().AsyncCallVoid(id, func(entity ec.Entity) {
			entity.DestroySelf()
		})
	}
}
