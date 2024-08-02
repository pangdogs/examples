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

package serv

import (
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/rpcpcsr"
)

// WorkService 工作服务
type WorkService struct {
	framework.ServiceInstance
}

func (serv *WorkService) InstallRPC(ctx service.Context) {
	// 安装RPC插件
	rpc.Install(ctx,
		rpc.With.Processors(
			rpcpcsr.NewServiceProcessor(nil),
			rpcpcsr.NewForwardProcessor(misc.Gate, gap.DefaultMsgCreator(), nil),
		),
	)
}

func (serv *WorkService) Built(ctx service.Context) {
	// 定义User实体原型
	serv.CreateEntityPT(misc.User).
		AddComponent(&comp.UserComp{}).
		AddComponent(&comp.CmdComp{}).
		Scope(ec.Scope_Global).
		Declare()
}

func (serv *WorkService) WakeUpUser(entId uid.Id) {
	// 创建User实体
	_, err := serv.CreateConcurrentEntity(misc.User).
		PersistId(entId).
		Spawn()
	if err != nil {
		log.Panic(serv, err)
	}
}
