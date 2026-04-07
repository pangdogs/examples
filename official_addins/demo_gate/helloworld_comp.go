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
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/generic"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

func (comp *HelloWorldComp) Awake() {
	session, ok := Gate.Require(service.Current(comp)).Get(comp.Entity().Id())
	if !ok {
		log.L(service.Current(comp)).Panic("session not found")
	}
	err := session.DataIO().Listen(comp.Entity(), generic.CastDelegateVoid2(comp.handleData))
	if err != nil {
		log.L(service.Current(comp)).Panic("listen data failed", zap.Error(err))
		return
	}
}

func (comp *HelloWorldComp) handleData(session gate.ISession, data []byte) {
	log.L(service.Current(comp)).Info("[receive]", zap.ByteString("data", data))

	echoData := pie.Reverse(data)
	err := session.DataIO().Send(echoData)
	if err != nil {
		log.L(service.Current(comp)).Panic("send data failed", zap.Error(err))
	}

	log.L(service.Current(comp)).Info("[send]", zap.ByteString("data", echoData))
}
