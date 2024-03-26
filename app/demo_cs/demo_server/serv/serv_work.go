package serv

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/service"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/processor"
)

type WorkService struct {
	framework.ServiceBehavior
}

func (serv *WorkService) Init(ctx service.Context) {
	core.CreateEntityPT(ctx).
		Prototype("user").
		AddComponent(&comp.UserComp{}).
		AddComponent(&comp.CmdComp{}).
		Declare()

	EntitiesPluginSelf.Install(ctx)
}

func (serv *WorkService) InstallRPC(ctx service.Context) {
	rpc.Install(ctx,
		rpc.With.Deliverers(
			processor.NewServiceDeliverer(),
			processor.NewForwardOutDeliverer(Gate),
		),
		rpc.With.Dispatchers(
			processor.NewServiceDispatcher(),
			processor.NewForwardInDispatcher(gap.DefaultMsgCreator()),
		),
	)
}
