package main

import (
	"fmt"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/logger"
	"kit.golaxy.org/plugins/registry"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[Demo, _Demo]("Demo组件")

// Demo Demo组件接口
type Demo interface{}

// _Demo Demo组件实现
type _Demo struct {
	ec.ComponentBehavior
}

// Update 组件更新
func (comp *_Demo) Update() {
	frame := runtime.Get(comp).GetFrame()

	if frame.GetCurFrames()%uint64(frame.GetTargetFPS()) == 0 {
		err := registry.Register(service.Get(comp), registry.Service{
			Name:    "demo",
			Version: "v0.1.0",
			Nodes: []registry.Node{
				{
					Id:      service.Get(comp).GetID().String(),
					Address: fmt.Sprintf("service:%s:%s", service.Get(comp).GetName(), service.Get(comp).GetID()),
				},
			},
		}, 3*time.Second)
		if err != nil {
			logger.Panic(service.Get(comp), err)
		}
	}
}

// Shut 组件停止
func (comp *_Demo) Shut() {
	err := registry.Deregister(service.Get(comp), registry.Service{
		Name:    "demo",
		Version: "v0.1.0",
		Nodes: []registry.Node{
			{
				Id:      service.Get(comp).GetID().String(),
				Address: fmt.Sprintf("service:%s:%s", service.Get(comp).GetName(), service.Get(comp).GetID()),
			},
		},
	})
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}
}
