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
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/rpc/callpath"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件实现
type HelloWorldComp struct {
	ec.ComponentBehavior
}

func (comp *HelloWorldComp) Start() {
	core.Await(runtime.Current(comp),
		core.TimeTickAsync(runtime.Current(comp), 3*time.Second),
	).Foreach(func(ctx runtime.Context, _ async.Result, _ ...any) {
		cp := callpath.CallPath{
			TargetKind: callpath.Entity,
			ExcludeSrc: true,
			Id:         entityId,
			Script:     "HelloWorldComp",
			Method:     "TestRPC",
		}

		n := rand.Uint32()
		err := RPC.Require(service.Current(comp)).OnewayRPC(Dsvc.Require(service.Current(comp)).NodeDetails().BalanceAddr, nil, cp, n)
		if err != nil {
			log.L(runtime.Current(comp)).Panic("oneway rpc failed", zap.Error(err))
		}

		log.L(runtime.Current(comp)).Info("[TestRPC] =>", zap.Uint32("n", n))
	})
}

func (comp *HelloWorldComp) TestRPC(n uint32) {
	log.L(runtime.Current(comp)).Info("=> [TestRPC]", zap.Uint32("n", n))
}
