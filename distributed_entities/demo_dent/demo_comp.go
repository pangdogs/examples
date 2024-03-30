package main

import (
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/framework/oc"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	oc.ComponentBehavior
}

func (comp *DemoComp) Start() {
	comp.Await(comp.TimeTick(5*time.Second)).Pipe(nil, func(_ runtime.Ret, _ ...any) {
		comp.GlobalBalanceOneWayRPC("DemoComp", "TestOnewayRPC", comp.GetServiceCtx().GetName(), comp.GetServiceCtx().GetId().String(), rand.Int31())
	})
}

func (comp *DemoComp) TestOnewayRPC(serv, id string, a int) {
	log.Infof(comp.GetRuntimeCtx(), "from: %s:%s => accept: %d", serv, id, a)
}
