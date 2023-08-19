package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/logger"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[any, DemoComp]("Demo组件")

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

// Awake 组件唤醒
func (comp *DemoComp) Awake() {
	logger.Infof(service.Current(comp), "I'm entity %q, comp %q Awake.", comp.GetEntity(), comp)
}

// Start 组件开始
func (comp *DemoComp) Start() {
	logger.Infof(service.Current(comp), "I'm entity %q, comp %q Start.", comp.GetEntity(), comp)
}

// Update 组件更新
func (comp *DemoComp) Update() {
	frame := runtime.Current(comp).GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		logger.Infof(service.Current(comp), "I'm entity %q, comp %q Update(%s).", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// LateUpdate 组件滞后更新
func (comp *DemoComp) LateUpdate() {
	ctx := runtime.Current(comp)
	frame := ctx.GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		logger.Infof(service.Current(comp), "I'm entity %q, comp %q LateUpdate(%s).", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// Shut 组件停止
func (comp *DemoComp) Shut() {
	logger.Infof(service.Current(comp), "I'm entity %q, comp %q Shut.", comp.GetEntity(), comp)
}
