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
	"git.golaxy.org/framework/plugins/broker"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	sub      broker.ISubscriber
	sequence int
}

// Start 组件开始
func (comp *DemoComp) Start() {
	log.Infof(service.Current(comp), "max payload: %d", broker.Using(service.Current(comp)).GetMaxPayload())

	sub, err := broker.Using(service.Current(comp)).Subscribe(context.Background(), "demo.>",
		broker.With.EventHandler(generic.MakeDelegateFunc1(func(e broker.IEvent) error {
			log.Infof(service.Current(comp), "receive=> pattern:%q, topic:%q, msg:%q", e.Pattern(), e.Topic(), string(e.Message()))
			return nil
		})))
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
	comp.sub = sub

	core.Await(runtime.Current(comp),
		core.TimeTick(runtime.Current(comp), time.Duration(rand.Int63n(5000))*time.Millisecond),
	).Pipe(nil, func(ctx runtime.Context, _ async.Ret, _ ...any) {
		topic := "demo.broker_test"
		msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)

		if err := broker.Using(service.Current(comp)).Publish(context.Background(), topic, []byte(msg)); err != nil {
			log.Panic(service.Current(comp), err)
		}

		log.Infof(service.Current(comp), "send=> topic:%q, msg:%q", topic, msg)
		comp.sequence++
	})
}

// Shut 组件结束
func (comp *DemoComp) Shut() {
	comp.sub.Unsubscribe()
}
