package main

import (
	"context"
	"encoding/json"
	"errors"
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
	service registry.Service
}

// Start 组件开始
func (comp *_Demo) Start() {
	w, err := registry.Watch(service.Get(comp), context.Background(), "demo")
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}

	comp.service = registry.Service{
		Name:    service.Get(comp).GetName(),
		Version: "v0.1.0",
		Nodes: []registry.Node{
			{
				Id:      service.Get(comp).GetId().String(),
				Address: fmt.Sprintf("service:%s:%s", service.Get(comp).GetName(), service.Get(comp).GetId()),
			},
		},
	}

	err = registry.Register(service.Get(comp), context.Background(), comp.service, 10*time.Second)
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}

	go func() {
		for {
			event, err := w.Next()
			if err != nil {
				if errors.Is(err, registry.ErrWatcherStopped) {
					logger.Info(service.Get(comp), "stop watching")
					return
				}
				logger.Panic(service.Get(comp), err)
			}

			eventData, _ := json.Marshal(event)
			logger.Infof(service.Get(comp), "receive event: %s", eventData)
		}
	}()
}

// Update 组件更新
func (comp *_Demo) Update() {
	frame := runtime.Get(comp).GetFrame()

	if frame.GetCurFrames()%uint64(150) == 0 {
		err := registry.Register(service.Get(comp), context.Background(), comp.service, 10*time.Second)
		if err != nil {
			logger.Panic(service.Get(comp), err)
		}
	}

	if frame.GetCurFrames()%uint64(300) == 0 {
		servces, err := registry.ListServices(service.Get(comp), context.Background())
		if err != nil {
			logger.Panic(service.Get(comp), err)
		}

		servicesData, _ := json.Marshal(servces)
		logger.Infof(service.Get(comp), "all services: %s", servicesData)
	}
}

// Shut 组件停止
func (comp *_Demo) Shut() {
	err := registry.Deregister(service.Get(comp), context.Background(), comp.service)
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}
}
