package main

import (
	"context"
	"errors"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/plugins/dsync"
	"git.golaxy.org/plugins/log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	mutex dsync.IDistMutex
}

// Update 组件更新
func (comp *DemoComp) Update() {
	if comp.mutex != nil {
		return
	}

	mutex := dsync.NewMutex(service.Current(comp), "demo_dsync_counter")
	if err := mutex.Lock(context.Background()); err != nil {
		log.Errorf(service.Current(comp), "lock failed: %s", err)
		return
	}
	comp.mutex = mutex

	log.Info(service.Current(comp), "lock")

	content, err := os.ReadFile("demo_dsync_counter.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic(service.Current(comp), err)
		}
	}

	n, _ := strconv.Atoi(string(content))
	n++

	log.Infof(service.Current(comp), "counter: %d", n)

	err = os.WriteFile("demo_dsync_counter.txt", []byte(strconv.Itoa(n)), os.ModePerm)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	core.Await(runtime.Current(comp),
		core.TimeAfter(context.Background(), time.Duration(rand.Int63n(1000))*time.Millisecond),
	).Any(runtime.Current(comp), func(ctx runtime.Context, _ runtime.Ret, _ ...any) {
		if comp.mutex == nil {
			return
		}
		comp.mutex.Unlock(context.Background())
		comp.mutex = nil

		log.Info(service.Current(comp), "unlock")
	})
}
