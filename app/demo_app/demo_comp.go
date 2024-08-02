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
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"time"
)

// DemoCompSelf Demo组件定义
var DemoCompSelf = define.Component[DemoComp]()

// DemoComp Demo组件
type DemoComp struct {
	framework.ComponentBehavior
}

func (comp *DemoComp) Start() {
	// 每隔5秒，以均衡模式，发送测试单程RPC
	comp.Await(comp.TimeTick(5*time.Second)).Pipe(nil, func(async.Ret, ...any) {
		comp.GlobalBalanceOneWayRPC(true, DemoCompSelf.Name, "TestOnewayRPC", rand.Int31())
	})
}

func (comp *DemoComp) TestOnewayRPC(r int) {
	log.Infof(comp, "entityId: %s, callChain: %+v => accept: %d",
		comp.GetEntity().GetId(), comp.GetRuntime().GetRPCStack().CallChain(), r)
}
