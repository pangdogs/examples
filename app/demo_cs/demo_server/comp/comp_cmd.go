package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/framework/oc"
	"git.golaxy.org/framework/plugins/log"
)

var CmdCompSelf = define.DefineComponent[CmdComp]()

type CmdComp struct {
	oc.ComponentBehavior
}

func (comp *CmdComp) Echo(text string) string {
	log.Infof(comp.GetRuntimeCtx(), text)
	return text
}
