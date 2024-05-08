package serv

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/uid"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/processor"
)

type WorkService struct {
	framework.ServiceGeneric
}

func (serv *WorkService) Instantiation() service.Context {
	return &WorkServiceInst{}
}

func (serv *WorkService) Built(ctx service.Context) {
	core.CreateEntityPT(ctx).
		Prototype("user").
		AddComponent(&comp.UserComp{}).
		AddComponent(&comp.CmdComp{}).
		Declare()
}

func (serv *WorkService) InstallRPC(ctx service.Context) {
	rpc.Install(ctx,
		rpc.With.Processors(
			processor.NewServiceProcessor(),
			processor.NewForwardProcessor(misc.Gate, gap.DefaultMsgCreator(), nil),
		),
	)
}

type WorkServiceInst struct {
	framework.ServiceInstance
}

func (inst *WorkServiceInst) CreateEntity(entId uid.Id) {
	rt := framework.CreateRuntime(inst).Spawn()

	runtime.Concurrent(rt).CallVoid(func(...any) {
		_, err := core.CreateEntity(rt).
			Prototype("user").
			Scope(ec.Scope_Global).
			PersistId(entId).
			Spawn()
		if err != nil {
			panic(err)
		}
	}).Wait(inst)
}
