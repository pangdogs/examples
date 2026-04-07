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
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/examples/app/demo_chat/server/comps"
	"git.golaxy.org/framework"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/rpcpcsr"
	"git.golaxy.org/framework/net/gap"
	"go.uber.org/zap"
)

// ChatService 聊天服务
type ChatService struct {
	framework.ServiceBehavior
}

func (s *ChatService) OnBuilt(svc framework.IService) {
	// 定义用户实体原型
	s.BuildEntityPT(consts.User).
		SetScope(ec.Scope_Global).
		AddComponent(&comps.ChatUserComp{}).
		Declare()
}

func (s *ChatService) InstallRPC(svc framework.IService) {
	// 安装RPC插件
	RPC.Install(s,
		rpc.With.Processors(
			rpcpcsr.NewServiceProcessor(nil, true),
			rpcpcsr.NewForwardProcessor(consts.Gate, gap.DefaultMsgCreator(), generic.CastDelegate2(rpcpcsr.DefaultValidateCliPermission), true),
		),
	)
}

func (s *ChatService) WakeUpUser(userId uid.Id) {
	// 创建用户实体
	user, err := s.BuildEntity(consts.User).
		SetPersistId(userId).
		New()
	if err != nil {
		s.L().Panic("create user failed", log.JSONRawStringer("user", user), zap.Error(err))
		return
	}
	s.L().Info("user created", log.JSONRawStringer("user", user))
}
