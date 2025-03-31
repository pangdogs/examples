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
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/framework/addins/dsvc"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/callpath"
	"git.golaxy.org/framework/utils/concurrent"
	"math/rand"
	"time"
)

type DstEntity struct {
	Id   uid.Id
	Addr string
}

var entities = concurrent.MakeLockedSlice[*DstEntity](0, 0)

// HelloWorldComp HelloWorld组件实现
type HelloWorldComp struct {
	ec.ComponentBehavior
}

func (comp *HelloWorldComp) Start() {
	entities.Append(&DstEntity{
		Id:   comp.GetId(),
		Addr: dsvc.Using(service.Current(comp)).GetNodeDetails().LocalAddr,
	})

	core.Await(runtime.Current(comp),
		core.TimeTickAsync(runtime.Current(comp), 3*time.Second),
	).Foreach(func(ctx runtime.Context, _ async.Ret, _ ...any) {
		var dstEntity *DstEntity

		entities.AutoRLock(func(es *[]*DstEntity) {
			if len(*es) <= 0 {
				return
			}
			dstEntity = (*es)[rand.Intn(len(*es))]
		})

		if dstEntity.Id == comp.GetId() {
			return
		}

		cp := callpath.CallPath{
			Category: callpath.Entity,
			Id:       dstEntity.Id,
			Script:   "HelloWorldComp",
			Method:   "TestRPC",
		}

		a := rand.Uint32()

		rv, err := rpc.Result1[int32](<-rpc.Using(service.Current(comp)).RPC(dstEntity.Addr, nil, cp, a)).Extract()
		if err != nil {
			log.Errorf(service.Current(comp), "send: %d, result: %v", a, err)
		} else {
			log.Infof(service.Current(comp), "send: %d, result: %d", a, rv)
		}
	})
}

func (comp *HelloWorldComp) Shut() {
	entities.DeleteOnce(func(entity *DstEntity) bool {
		return entity.Id == comp.GetId()
	})
}

func (comp *HelloWorldComp) TestRPC(a uint32) int32 {
	n := rand.Int31()
	log.Infof(service.Current(comp), "accept: %d, return: %d", a, n)
	return n
}
