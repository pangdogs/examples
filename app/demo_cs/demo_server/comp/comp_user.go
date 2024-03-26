package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/framework/oc"
)

var UserCompSelf = define.DefineComponent[UserComp]()

type UserComp struct {
	oc.ComponentBehavior
}

func (c *UserComp) Dispose() {
	runtime.Concurrent(c).GetCancelFunc()()
}
