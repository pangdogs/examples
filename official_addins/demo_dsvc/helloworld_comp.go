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
	"math/rand"
	"time"

	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gap/variant"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	// 监听消息
	Dsvc.Require(service.Current(comp)).Listen(comp.Entity(), generic.CastDelegateVoid2(
		func(topic string, mp gap.MsgPacket) {
			log.L(service.Current(comp)).Info("[receive]",
				zap.String("topic", topic),
				log.JSON("msg-packet", mp))
		},
	))

	// 定时发送消息
	core.Await(comp.Entity(), core.TimeTickAsync(comp.Entity(), time.Second)).
		Foreach(func(ctx runtime.Context, ret async.Result, _ ...any) {
			details := Dsvc.Require(service.Current(ctx)).NodeDetails()

			m, _ := variant.NewMapFromGoMap(map[string]int{
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
			})
			arr, _ := variant.NewArray([]int{rand.Int(), rand.Int(), rand.Int()})

			msg := &MsgHelloWorld{
				Int:    rand.Int(),
				Double: rand.Float64(),
				Str:    ksuid.New().String(),
				Map:    m,
				Array:  arr,
			}

			// 广播消息
			err := Dsvc.Require(service.Current(ctx)).Send(details.BroadcastAddr, msg)
			if err != nil {
				log.L(service.Current(ctx)).Error("send failed", zap.Error(err))
			}

			log.L(service.Current(comp)).Info("[send]",
				zap.String("addr", details.BroadcastAddr),
				log.JSON("msg", msg))
		})
}
