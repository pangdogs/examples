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
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/log"
	"math/rand"
	"time"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	framework.ComponentBehavior
}

func (comp *HelloWorldComp) Start() {
	// 每隔3秒，测试广播单程RPC
	comp.Await(comp.TimeTickAsync(3 * time.Second)).Foreach(func(async.Ret, ...any) {
		comp.GlobalBroadcastOnewayRPC(true, comp.GetName(), "TestOnewayRPC", rand.Int31())
	})

	// 10秒后销毁实体
	comp.Await(comp.TimeAfterAsync(10 * time.Second)).Foreach(func(async.Ret, ...any) {
		comp.GetEntity().DestroySelf()
	})
}

func (comp *HelloWorldComp) TestOnewayRPC(a int) {
	log.Infof(comp, "callChain: %+v => accept: %d", comp.GetRuntime().GetRPCStack().CallChain(), a)
}
