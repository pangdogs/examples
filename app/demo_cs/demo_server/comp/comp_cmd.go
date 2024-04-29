package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/framework/fwec"
	"git.golaxy.org/framework/plugins/log"
)

var CmdCompSelf = define.Component[CmdComp]()

type CmdComp struct {
	fwec.ComponentBehavior
}

func (comp *CmdComp) Echo(text string) string {
	log.Infof(comp.GetRuntime().Ctx, text)
	return text
}
