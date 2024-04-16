package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/framework/dc"
)

var UserCompSelf = define.DefineComponent[UserComp]()

type UserComp struct {
	dc.ComponentBehavior
}

func (c *UserComp) Dispose() {
	runtime.Concurrent(c).GetCancelFunc()()
}
