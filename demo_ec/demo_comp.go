package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/plugins/logger"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[Demo, _Demo]("Demo组件")

// Demo Demo组件接口
type Demo interface{}

// _Demo Demo组件实现
type _Demo struct {
	ec.ComponentBehavior
}

// Awake 组件唤醒
func (comp *_Demo) Awake() {
	logger.Infof(runtime.Get(comp), "I'm entity %q, comp %q Awake.", comp.GetEntity(), comp)
}

// Start 组件开始
func (comp *_Demo) Start() {
	logger.Infof(runtime.Get(comp), "I'm entity %q, comp %q Start.", comp.GetEntity(), comp)
}

// Update 组件更新
func (comp *_Demo) Update() {
	ctx := runtime.Get(comp)
	frame := ctx.GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		logger.Infof(runtime.Get(comp), "I'm entity %q, comp %q Update(%s).", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// LateUpdate 组件滞后更新
func (comp *_Demo) LateUpdate() {
	ctx := runtime.Get(comp)
	frame := ctx.GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		logger.Infof(runtime.Get(comp), "I'm entity %q, comp %q LateUpdate(%s).", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// Shut 组件停止
func (comp *_Demo) Shut() {
	logger.Infof(runtime.Get(comp), "I'm entity %q, comp %q Shut.", comp.GetEntity(), comp)
}
