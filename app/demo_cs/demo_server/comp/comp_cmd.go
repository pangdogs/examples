package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/framework/dc"
	"git.golaxy.org/framework/plugins/log"
)

var CmdCompSelf = define.DefineComponent[CmdComp]()

type CmdComp struct {
	dc.ComponentBehavior
}

func (comp *CmdComp) Echo(text string) string {
	log.Infof(comp.GetRuntimeCtx(), text)
	return text
}
