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
	"fmt"
	"time"

	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/service"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/discovery"
	"git.golaxy.org/framework/addins/log"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
	service *discovery.Service
	reg     discovery.IRegistration
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	w, err := Discovery.Require(service.Current(comp)).WatchEvent(comp.Entity(), service.Current(comp).Name())
	if err != nil {
		log.L(service.Current(comp)).Panic("watch event failed", zap.Error(err))
	}

	comp.service = &discovery.Service{
		Name: service.Current(comp).Name(),
		Nodes: []discovery.Node{
			{
				Id:      service.Current(comp).Id(),
				Address: fmt.Sprintf("service:%s:%s", service.Current(comp).Name(), service.Current(comp).Id()),
			},
		},
	}

	reg, err := Discovery.Require(service.Current(comp)).RegisterNode(comp.Entity(), comp.service.Name, &comp.service.Nodes[0], 3*time.Second)
	if err != nil {
		log.L(service.Current(comp)).Panic("register node failed", zap.Error(err))
	}
	_, err = reg.KeepAliveContinuous(comp.Entity())
	if err != nil {
		log.L(service.Current(comp)).Panic("keep alive continuous failed", zap.Error(err))
	}
	comp.reg = reg

	go func() {
		for e := range w {
			if e.Error != nil {
				log.L(service.Current(comp)).Error("watch event failed", zap.Error(e.Error))
			}
			log.L(service.Current(comp)).Info("[event]", log.JSON("event", e))
		}
	}()
}

// Shut 组件停止
func (comp *HelloWorldComp) Shut() {
	err := comp.reg.Deregister(context.Background())
	if err != nil {
		log.L(service.Current(comp)).Panic("deregister node failed", zap.Error(err))
	}
}
