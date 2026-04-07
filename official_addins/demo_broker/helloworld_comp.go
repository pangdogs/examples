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
	"math/rand"
	"time"

	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/broker"
	"git.golaxy.org/framework/addins/log"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
	sequence int
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	log.L(service.Current(comp)).Info("starting...",
		zap.String("entity_id", comp.Entity().Id().String()),
		zap.Int64("max_payload", Broker.Require(service.Current(comp)).MaxPayload()))

	_, err := Broker.Require(service.Current(comp)).SubscribeHandler(comp.Entity(), "helloworld.>", "",
		generic.CastDelegateVoid1(func(e broker.Event) {
			log.L(service.Current(comp)).Info("[receive]",
				zap.String("pattern", e.Pattern),
				zap.String("topic", e.Topic),
				zap.ByteString("msg", e.Message))
		}),
	)
	if err != nil {
		log.L(service.Current(comp)).Panic("subscribe error", zap.Error(err))
	}

	core.Await(comp.Entity(),
		core.TimeTickAsync(comp.Entity(), time.Duration(rand.Int63n(3000))*time.Millisecond),
	).Foreach(func(ctx runtime.Context, _ async.Result, _ ...any) {
		topic := "helloworld.testing"
		msg := fmt.Sprintf("%s/%d", comp.Entity().Id(), comp.sequence)

		if err := Broker.Require(service.Current(comp)).Publish(context.Background(), topic, []byte(msg)); err != nil {
			log.L(service.Current(comp)).Panic("publish error", zap.Error(err))
		}

		log.L(service.Current(comp)).Info("[publish]",
			zap.String("topic", topic),
			zap.String("msg", msg))

		comp.sequence++
	})
}

// Shut 组件结束
func (comp *HelloWorldComp) Shut() {
	log.L(service.Current(comp)).Info("shutting...", zap.String("entity_id", comp.Entity().Id().String()))
}
