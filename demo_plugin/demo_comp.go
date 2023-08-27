package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[any, DemoComp]("demo组件")

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *DemoComp) Start() {
	// 调用demo插件
	defineDemoPlugin.Fetch(service.Current(comp)).HelloWorld()

	// 停止运行时
	runtime.Current(comp).GetCancelFunc()()
}
