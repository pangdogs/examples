package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/log"
	"kit.golaxy.org/plugins/registry"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	service registry.Service
}

// Start 组件开始
func (comp *DemoComp) Start() {
	w, err := registry.Watch(service.Current(comp), context.Background(), service.Current(comp).GetName())
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	comp.service = registry.Service{
		Name:    service.Current(comp).GetName(),
		Version: "v0.1.0",
		Nodes: []registry.Node{
			{
				Id:      service.Current(comp).GetId().String(),
				Address: fmt.Sprintf("service:%s:%s", service.Current(comp).GetName(), service.Current(comp).GetId()),
			},
		},
	}

	err = registry.Register(service.Current(comp), context.Background(), comp.service, 10*time.Second)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	go func() {
		for {
			event, err := w.Next()
			if err != nil {
				if errors.Is(err, registry.ErrStoppedWatching) {
					log.Info(service.Current(comp), "stop watching")
					return
				}
				log.Panic(service.Current(comp), err)
			}

			eventData, _ := json.Marshal(event)
			log.Infof(service.Current(comp), "receive event: %s", eventData)
		}
	}()
}

// Update 组件更新
func (comp *DemoComp) Update() {
	frame := runtime.Current(comp).GetFrame()

	if frame.GetCurFrames()%uint64(150) == 0 {
		err := registry.Register(service.Current(comp), context.Background(), comp.service, 10*time.Second)
		if err != nil {
			log.Panic(service.Current(comp), err)
		}
	}

	if frame.GetCurFrames()%uint64(300) == 0 {
		servces, err := registry.ListServices(service.Current(comp), context.Background())
		if err != nil {
			log.Panic(service.Current(comp), err)
		}

		servicesData, _ := json.Marshal(servces)
		log.Infof(service.Current(comp), "all services: %s", servicesData)
	}
}

// Shut 组件停止
func (comp *DemoComp) Shut() {
	err := registry.Deregister(service.Current(comp), context.Background(), comp.service)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
}
