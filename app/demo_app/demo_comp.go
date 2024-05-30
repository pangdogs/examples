package main

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"time"
)

// DemoCompSelf Demo组件定义
var DemoCompSelf = define.Component[DemoComp]()

// DemoComp Demo组件
type DemoComp struct {
	framework.ComponentBehavior
}

func (comp *DemoComp) Start() {
	// 每隔5秒，以均衡模式，发送测试单程RPC
	comp.Await(comp.TimeTick(5*time.Second)).Pipe(nil, func(async.Ret, ...any) {
		comp.GlobalBalanceOneWayRPC(DemoCompSelf.Name, "TestOnewayRPC", rand.Int31())
	})
}

func (comp *DemoComp) TestOnewayRPC(r int) {
	log.Infof(comp, "entityId: %s, callChain: %+v => accept: %d",
		comp.GetEntity().GetId(), comp.GetRuntime().GetRPCStack().CallChain(), r)
}
