package main

import (
	"context"
	"errors"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/dsync"
	"kit.golaxy.org/plugins/logger"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[any, DemoComp]("Demo组件")

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	mutex dsync.DMutex
}

// Update 组件更新
func (comp *DemoComp) Update() {
	if comp.mutex != nil {
		return
	}

	mutex := dsync.NewDMutex(service.Get(comp), "demo_dsync_counter", dsync.Option{}.Tries(64))
	if err := mutex.Lock(context.Background()); err != nil {
		logger.Errorf(service.Get(comp), "lock failed: %s", err)
		return
	}
	comp.mutex = mutex

	logger.Info(service.Get(comp), "lock")

	content, err := os.ReadFile("demo_dsync_counter.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Panic(service.Get(comp), err)
		}
	}

	n, _ := strconv.Atoi(string(content))
	n++

	logger.Infof(service.Get(comp), "counter: %d", n)

	err = os.WriteFile("demo_dsync_counter.txt", []byte(strconv.Itoa(n)), os.ModePerm)
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}

	golaxy.Await(comp, golaxy.AsyncTimeAfter(context.Background(), time.Duration(rand.Int63n(1000))*time.Millisecond),
		func(ctx runtime.Context, ret runtime.Ret) {
			if comp.mutex == nil {
				return
			}
			comp.mutex.Unlock(context.Background())
			comp.mutex = nil

			logger.Info(service.Get(comp), "unlock")
		})
}
