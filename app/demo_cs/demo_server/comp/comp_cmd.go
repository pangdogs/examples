package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/log"
)

var CmdCompSelf = define.Component[CmdComp]()

type CmdComp struct {
	framework.ComponentBehavior
}

func (comp *CmdComp) Echo(text string) string {
	log.Infof(comp, "text:%s, callChain:%+v", text, comp.GetRuntime().GetRPCStack().CallChain())
	return text
}
