package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework/dc"
)

var UserCompSelf = define.DefineComponent[UserComp]()

type UserComp struct {
	dc.ComponentBehavior
}

func (c *UserComp) Dispose() {
	if c.GetServiceCtx().GetName() == misc.Gate {
		c.BroadcastOneWayRPC(misc.Work, "", "DestroySelf")
	}
	runtime.Concurrent(c).GetCancelFunc()()
}
