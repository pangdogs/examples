package main

import (
	"context"
	"github.com/golaxy-kit/golaxy/define"
	"github.com/golaxy-kit/golaxy/ec"
	"github.com/golaxy-kit/golaxy/service"
	"github.com/golaxy-kit/plugins/registry"
	"time"
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

}

// Update 组件更新
func (comp *_DemoComp) Update() {
	// 服务上下文
	serviceCtx := service.Get(comp)

	// 注册服务
	registry.Context(serviceCtx).Register(
		context.Background(),
		registry.Service{
			Name:    "demo",
			Version: "1.0.0",
			Nodes: []registry.Node{
				{
					Id:       "1",
					Address:  "",
					Metadata: nil,
				},
			},
		},
		registry.RegisterOption.TTL(10*time.Second))
}
