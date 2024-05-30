package main

import (
	"git.golaxy.org/core/service"
	"git.golaxy.org/framework"
)

// DemoService Demo服务实例
type DemoService struct {
	framework.ServiceInstance
}

func (serv *DemoService) Built(ctx service.Context) {
	// 声明实体原型
	serv.CreateEntityPT("demo").
		AddComponent(&DemoComp{}).
		Declare()
}

func (serv *DemoService) Started(ctx service.Context) {
	// 创建实体
	serv.CreateConcurrentEntity("demo").
		PersistId("1").
		Spawn()
}
