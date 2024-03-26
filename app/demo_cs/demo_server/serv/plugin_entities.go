package serv

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/uid"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/log"
)

var (
	EntitiesPluginSelf = define.DefineServicePlugin(newEntities)
)

type IEntities any

func newEntities(...any) IEntities {
	return &_Entities{}
}

type _Entities struct {
	servCtx service.Context
}

func (e *_Entities) InitSP(ctx service.Context) {
	log.Infof(ctx, "init plugin %q", EntitiesPluginSelf.Name)
	e.servCtx = ctx
}

func (e *_Entities) ShutSP(ctx service.Context) {
	log.Infof(ctx, "shut plugin %q", EntitiesPluginSelf.Name)
}

func (e *_Entities) CreateEntity(entId uid.Id) {
	rt := framework.CreateRuntime(e.servCtx).Spawn()

	runtime.Concurrent(rt).CallVoid(func(...any) {
		_, err := core.CreateEntity(rt).
			Prototype("user").
			Scope(ec.Scope_Global).
			PersistId(entId).
			Spawn()
		if err != nil {
			panic(err)
		}
	}).Wait(e.servCtx)
}
