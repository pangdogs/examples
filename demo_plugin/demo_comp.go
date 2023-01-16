package main

import (
	"github.com/golaxy-kit/golaxy/define"
	"github.com/golaxy-kit/golaxy/ec"
	"github.com/golaxy-kit/golaxy/service"
)

func init() {
	// 注册Demo组件
	DemoCompPt.Register(_DemoComp{}, "demo组件")
}

// DemoCompPt 定义Demo组件原型
var DemoCompPt = define.DefineComponentInterface[DemoComp]().ComponentInterface()

// DemoComp Demo组件接口
type DemoComp interface{}

// _DemoComp Demo组件实现类
type _DemoComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *_DemoComp) Start() {
	DemoPlugin.Context(service.Get(comp)).Test()
}
