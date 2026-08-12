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
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	framework.ComponentBehavior
}

func (comp *HelloWorldComp) Start() {
	// 每隔3秒，测试广播单程RPC
	comp.scheduleBroadcast(3 * time.Second)

	// 10秒后销毁实体
	core.ContinueOnVoid(comp,
		core.After(comp.AsyncScope().Context(), 10*time.Second),
		func(_ runtime.Context, ret async.Result, _ ...any) {
			if ret.Error == nil {
				comp.Entity().Destroy()
			}
		})
}

func (comp *HelloWorldComp) scheduleBroadcast(interval time.Duration) {
	core.ContinueOnVoid(comp,
		core.After(comp.AsyncScope().Context(), interval),
		func(_ runtime.Context, ret async.Result, _ ...any) {
			if ret.Error != nil {
				return
			}

			n := rand.Int31()
			if err := comp.GlobalBroadcastOnewayRPC(true, comp.Name(), "TestOnewayRPC", n); err != nil {
				comp.L().Panic("TestOnewayRPC failed", zap.Error(err))
			}
			comp.L().Info("[TestOnewayRPC] =>",
				zap.Any("callChain", comp.Runtime().RPCStack().CallChain()),
				zap.Int32("n", n))
			comp.scheduleBroadcast(interval)
		})
}

func (comp *HelloWorldComp) TestOnewayRPC(n int) {
	comp.L().Info("=> [TestOnewayRPC]",
		zap.Any("callChain", comp.Runtime().RPCStack().CallChain()),
		zap.Int("n", n))
}
