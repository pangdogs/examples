/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/framework/plugins/discovery"
	"git.golaxy.org/framework/plugins/log"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	service *discovery.Service
}

// Start 组件开始
func (comp *DemoComp) Start() {
	w, err := discovery.Using(service.Current(comp)).Watch(context.Background(), service.Current(comp).GetName())
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	comp.service = &discovery.Service{
		Name: service.Current(comp).GetName(),
		Nodes: []discovery.Node{
			{
				Id:      service.Current(comp).GetId(),
				Address: fmt.Sprintf("service:%s:%s", service.Current(comp).GetName(), service.Current(comp).GetId()),
			},
		},
	}

	err = discovery.Using(service.Current(comp)).Register(context.Background(), comp.service, 10*time.Second)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	go func() {
		for {
			event, err := w.Next()
			if err != nil {
				if errors.Is(err, discovery.ErrTerminated) {
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

	if frame.GetCurFrames()%int64(150) == 0 {
		err := discovery.Using(service.Current(comp)).Register(context.Background(), comp.service, 10*time.Second)
		if err != nil {
			log.Panic(service.Current(comp), err)
		}
	}

	if frame.GetCurFrames()%int64(300) == 0 {
		servces, err := discovery.Using(service.Current(comp)).ListServices(context.Background())
		if err != nil {
			log.Panic(service.Current(comp), err)
		}

		servicesData, _ := json.Marshal(servces)
		log.Infof(service.Current(comp), "all services: %s", servicesData)
	}
}

// Shut 组件停止
func (comp *DemoComp) Shut() {
	err := discovery.Using(service.Current(comp)).Deregister(context.Background(), comp.service)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
}
