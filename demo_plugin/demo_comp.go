package main

import (
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *DemoComp) Start() {
	// 调用demo插件
	Using(service.Current(comp)).HelloWorld()
	// 停止运行时
	runtime.Current(comp).GetCancelFunc()()
}
