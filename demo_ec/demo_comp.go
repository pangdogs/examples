package main

import (
	"fmt"
	"github.com/golaxy-kit/golaxy/define"
	"github.com/golaxy-kit/golaxy/ec"
)

// DemoComp 定义Demo组件
var DemoComp = define.DefineComponent[Demo, _Demo]()

// Demo Demo组件接口
type Demo interface{}

// _Demo Demo组件
type _Demo struct {
	ec.ComponentBehavior
	count int
}

// Awake 组件唤醒
func (comp *_Demo) Awake() {
	fmt.Printf("I'm entity %s, %s Awake.\n", comp.GetEntity(), comp)
}

// Start 组件开始
func (comp *_Demo) Start() {
	fmt.Printf("I'm entity %s, %s Start.\n", comp.GetEntity(), comp)
}

// Update 组件更新
func (comp *_Demo) Update() {
	if comp.count%30 == 0 {
		fmt.Printf("I'm entity %s, %s Update(%d).\n", comp.GetEntity(), comp, comp.count)
	}
}

// LateUpdate 组件滞后更新
func (comp *_Demo) LateUpdate() {
	if comp.count%30 == 0 {
		fmt.Printf("I'm entity %s, %s LateUpdate(%d).\n", comp.GetEntity(), comp, comp.count)
	}
	comp.count++
}

// Shut 组件停止
func (comp *_Demo) Shut() {
	fmt.Printf("I'm entity %s, %s Shut.\n", comp.GetEntity(), comp)
}
