package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
)

var UserCompSelf = define.Component[UserComp]()

type UserComp struct {
	framework.ComponentBehavior
}

func (c *UserComp) Dispose() {
	if c.GetService().Ctx.GetName() == misc.Gate {
		c.BroadcastOneWayRPC(misc.Work, "", "DestroySelf")
	}
	runtime.Concurrent(c).GetCancelFunc()()
}
