package main

import (
	"fmt"
	"github.com/golaxy-kit/golaxy/define"
	"github.com/golaxy-kit/golaxy/ec"
	"github.com/golaxy-kit/golaxy/runtime"
)

// DemoComp 定义Demo组件
var DemoComp = define.DefineComponent[Demo, _Demo]("Demo组件，在组件生命周期回调函数中，打印一些信息。")

// Demo Demo组件接口
type Demo interface{}

// _Demo Demo组件
type _Demo struct {
	ec.ComponentBehavior
}

// Awake 组件唤醒
func (comp *_Demo) Awake() {
	fmt.Printf("I'm entity %s, comp %s Awake.\n", comp.GetEntity(), comp)
}

// Start 组件开始
func (comp *_Demo) Start() {
	fmt.Printf("I'm entity %s, comp %s Start.\n", comp.GetEntity(), comp)
}

// Update 组件更新
func (comp *_Demo) Update() {
	ctx := runtime.Get(comp)
	frame := ctx.GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		fmt.Printf("I'm entity %s, comp %s Update(%f).\n", comp.GetEntity(), comp, frame.GetRunningElapseTime().Seconds())
	}
}

// LateUpdate 组件滞后更新
func (comp *_Demo) LateUpdate() {
	ctx := runtime.Get(comp)
	frame := ctx.GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		fmt.Printf("I'm entity %s, comp %s LateUpdate(%f).\n", comp.GetEntity(), comp, frame.GetRunningElapseTime().Seconds())
	}
}

// Shut 组件停止
func (comp *_Demo) Shut() {
	fmt.Printf("I'm entity %s, comp %s Shut.\n", comp.GetEntity(), comp)
}
