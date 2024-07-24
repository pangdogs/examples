package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/rpc/rpcutil"
)

var UserCompSelf = define.Component[UserComp]()

type UserComp struct {
	framework.ComponentBehavior
}

func (c *UserComp) Dispose() {
	if c.GetService().GetName() == misc.Gate {
		c.RPC(misc.Work, rpcutil.NoComp, "DestroySelf")
	}
}
