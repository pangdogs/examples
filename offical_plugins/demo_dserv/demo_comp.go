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
	"encoding/json"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework/net/gap/variant"
	"git.golaxy.org/framework/plugins/dserv"
	"git.golaxy.org/framework/plugins/log"
	"github.com/segmentio/ksuid"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

func (comp *DemoComp) Start() {
	core.Await(runtime.Current(comp), core.TimeTick(runtime.Current(comp), time.Second)).
		Pipe(runtime.Current(comp), func(ctx runtime.Context, ret async.Ret, _ ...any) {
			addr := dserv.Using(service.Current(ctx)).GetNodeDetails()

			vmap, err := variant.MakeReadonlyMapFromGoMap(map[string]int{
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

			msg := &MsgDemo{
				Int:    rand.Int(),
				Double: rand.Float64(),
				Str:    ksuid.New().String(),
				Map:    vmap,
				Array:  arr,
			}

			// 广播消息
			err = dserv.Using(service.Current(ctx)).SendMsg(addr.BroadcastAddr, msg)
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			data, _ := json.Marshal(msg)
			log.Infof(service.Current(ctx), "send => topic:%q, msg:%s", addr.BroadcastAddr, data)
		})
}
