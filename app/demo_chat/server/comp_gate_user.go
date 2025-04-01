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
	"git.golaxy.org/examples/app/demo_chat/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/router"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/rpcutil"
)

type GateUserComp struct {
	framework.ComponentBehavior
}

func (c *GateUserComp) Start() {
	session := c.GetEntity().GetMeta().Value("session").(gate.ISession)

	mapping, err := router.Using(c.GetService()).Mapping(c.GetId(), session.GetId())
	if err != nil {
		log.Panicf(c, "mapping gate user %s to session %s failed, %s", c.GetId(), session.GetId(), err)
	}

	err = rpc.ResultVoid(<-rpcutil.ProxyService(c.GetService(), misc.Chat).BalanceRPC(rpcutil.NoAddIn, "WakeUpUser", c.GetId())).Extract()
	if err != nil {
		log.Panicf(c, "wakeup chat user %s failed, %s", c.GetId(), err)
	}

	go func() {
		<-mapping.Done()
		<-c.GetRuntime().Terminate()
	}()
}

func (c *GateUserComp) Shut() {
	<-c.RPC(misc.Chat, "ChatUserComp", "DestroySelf")
}

func (c *GateUserComp) Dispose() {
	log.Infof(c, "gate user %s destroyed", c.GetId())
}
