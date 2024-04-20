package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework/fwec"
)

var UserCompSelf = define.DefineComponent[UserComp]()

type UserComp struct {
	fwec.ComponentBehavior
}

func (c *UserComp) Dispose() {
	if c.GetService().Ctx.GetName() == misc.Gate {
		c.BroadcastOneWayRPC(misc.Work, "", "DestroySelf")
	}
	runtime.Concurrent(c).GetCancelFunc()()
}
