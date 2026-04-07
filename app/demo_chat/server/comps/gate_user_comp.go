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

package comps

import (
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/framework"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/rpc"
	"go.uber.org/zap"
)

type GateUserComp struct {
	framework.ComponentBehavior
	chatChannel *GateChatChannelComp
}

func (c *GateUserComp) Start() {
	mapping, err := Router.Require(c.Service()).Map(c.Id(), c.Id())
	if err != nil {
		c.L().Panic("mapping failed", zap.Any("user", c.Entity()), zap.Error(err))
	}

	err = rpc.ResultVoid(rpc.ProxyService(c).BalanceRPC(consts.Chat, "", "WakeUpUser", c.Id())).Error
	if err != nil {
		c.L().Panic("RPC::WakeUpUser failed", zap.Any("user", c.Entity()), zap.Error(err))
	}

	go func() {
		<-mapping.Unmapped().Done()
		<-c.Runtime().Terminate().Done()
	}()
}

func (c *GateUserComp) Shut() {
	err := rpc.ResultVoid(c.RPC(consts.Chat, "ChatUserComp", "Destroy")).Error
	if err != nil {
		c.L().Error("RPC::ChatUserComp.Destroy failed", zap.Any("user", c.Entity()), zap.Error(err))
		return
	}
}

func (c *GateUserComp) Dispose() {
	c.L().Info("user disposed")
}
