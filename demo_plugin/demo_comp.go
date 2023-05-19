package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[DemoComp, _DemoComp]("demo组件")

// DemoComp Demo组件接口
type DemoComp interface{}

// _DemoComp Demo组件实现
type _DemoComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *_DemoComp) Start() {
	// 调用demo插件
	defineDemoPlugin.Get(service.Get(comp)).HelloWorld()

	// 停止运行时
	runtime.Get(comp).GetCancelFunc()()
}
