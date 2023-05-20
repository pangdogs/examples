package main

import (
	"context"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/dsync"
	"kit.golaxy.org/plugins/logger"
	"math/rand"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[any, DemoComp]("Demo组件")

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

// Update 组件更新
func (comp *DemoComp) Update() {
	mutex := dsync.NewDMutex(service.Get(comp), "demo_dsync_count")
	if err := mutex.Lock(context.Background()); err != nil {
		logger.Panic(service.Get(comp), err)
	}
	defer mutex.Unlock(context.Background())

	sleepTime := time.Duration(rand.Intn(2000)) * time.Millisecond
	time.Sleep(sleepTime)

	logger.Panicf(service.Get(comp), "sleep: %f", sleepTime.Seconds())
}
