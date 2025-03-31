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
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/addins/dsvc"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gap/variant"
	"github.com/segmentio/ksuid"
	"math/rand"
	"time"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	// 监听消息
	dsvc.Using(service.Current(comp)).WatchMsg(context.Background(), generic.CastDelegate2(
		func(topic string, mp gap.MsgPacket) error {
			data, _ := json.Marshal(mp)
			log.Infof(service.Current(comp), "=>[receive] topic:%q, msg-packet:%s", topic, data)
			return nil
		},
	))

	// 定时发送消息
	core.Await(runtime.Current(comp), core.TimeTickAsync(runtime.Current(comp), time.Second)).
		Foreach(func(ctx runtime.Context, ret async.Ret, _ ...any) {
			details := dsvc.Using(service.Current(ctx)).GetNodeDetails()

			m, err := variant.MakeReadonlyMapFromGoMap(map[string]int{
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
			})
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			arr, err := variant.MakeReadonlyArray([]int{rand.Int(), rand.Int(), rand.Int()})
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			msg := &MsgHelloWorld{
				Int:    rand.Int(),
				Double: rand.Float64(),
				Str:    ksuid.New().String(),
				Map:    m,
				Array:  arr,
			}

			// 广播消息
			err = dsvc.Using(service.Current(ctx)).SendMsg(details.BroadcastAddr, msg)
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			data, _ := json.Marshal(msg)
			log.Infof(service.Current(ctx), "[send]=> topic:%q, msg:%s", details.BroadcastAddr, data)
		})
}
