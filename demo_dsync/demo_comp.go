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
var defineDemoComp = define.DefineComponent[Demo, _Demo]("Demo组件")

// Demo Demo组件接口
type Demo interface{}

// _Demo Demo组件实现
type _Demo struct {
	ec.ComponentBehavior
}

// Update 组件更新
func (comp *_Demo) Update() {
	mutex := dsync.NewDMutex(service.Get(comp), "demo_dsync_count")
	if err := mutex.Lock(context.Background()); err != nil {
		logger.Panic(service.Get(comp), err)
	}
	defer mutex.Unlock(context.Background())

	sleepTime := time.Duration(rand.Intn(2000)) * time.Millisecond
	time.Sleep(sleepTime)

	logger.Panicf(service.Get(comp), "sleep: %f", sleepTime.Seconds())
}
