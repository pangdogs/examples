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
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/addins/broker"
	"git.golaxy.org/framework/addins/log"
	"math/rand"
	"time"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
	sub      broker.ISubscriber
	sequence int
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	log.Infof(service.Current(comp), "max payload: %d", broker.Using(service.Current(comp)).GetMaxPayload())

	sub, err := broker.Using(service.Current(comp)).Subscribe(context.Background(), "helloworld.>",
		broker.With.EventHandler(generic.CastDelegate1(func(e broker.IEvent) error {
			log.Infof(service.Current(comp), "=>[receive] pattern:%s, topic:%s, msg:%s", e.Pattern(), e.Topic(), string(e.Message()))
			return nil
		})))
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
	comp.sub = sub

	core.Await(runtime.Current(comp),
		core.TimeTickAsync(runtime.Current(comp), time.Duration(rand.Int63n(3000))*time.Millisecond),
	).Foreach(func(ctx runtime.Context, _ async.Ret, _ ...any) {
		topic := "helloworld.testing"
		msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)

		if err := broker.Using(service.Current(comp)).Publish(context.Background(), topic, []byte(msg)); err != nil {
			log.Panic(service.Current(comp), err)
		}

		log.Infof(service.Current(comp), "[send]=> topic:%s, msg:%s", topic, msg)
		comp.sequence++
	})
}

// Shut 组件结束
func (comp *HelloWorldComp) Shut() {
	comp.sub.Unsubscribe()
}
