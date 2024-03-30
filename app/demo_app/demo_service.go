package main

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/framework"
)

// DemoService Demo服务
type DemoService struct {
	framework.ServiceGeneric
}

func (serv *DemoService) Init(ctx service.Context) {
	// 声明实体原型
	core.CreateEntityPT(ctx).
		Prototype("demo").
		AddComponent(&DemoComp{}).
		Declare()
}

func (serv *DemoService) Started(ctx service.Context) {
	// 创建运行时
	rt := framework.CreateRuntime(ctx).Spawn()

	runtime.Concurrent(rt).CallVoid(func(...any) {
		// 创建实体
		core.CreateEntity(rt).
			Prototype("demo").
			Scope(ec.Scope_Global).
			PersistId("xxxxxx").
			Spawn()
	})
}
