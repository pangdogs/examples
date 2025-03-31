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
	"fmt"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"sync"
	"time"
)

var (
	textQueue []string
	textMutex sync.RWMutex
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
	session gate.ISession
	pos     int
}

func (comp *HelloWorldComp) Awake() {
	comp.session = comp.GetEntity().GetMeta().Value("session").(gate.ISession)
	comp.pos = len(textQueue)
}

func (comp *HelloWorldComp) Start() {
	textMutex.RLock()
	defer textMutex.RUnlock()

	err := comp.session.GetSettings().RecvDataHandler(generic.CastDelegate2(comp.onRecvData)).Change()
	if err != nil {
		log.Panic(runtime.Current(comp), err)
	}

	core.Await(runtime.Current(comp),
		core.TimeTickAsync(runtime.Current(comp), time.Second),
	).Foreach(func(ctx runtime.Context, ret async.Ret, _ ...any) {
		textMutex.RLock()
		defer textMutex.RUnlock()

		for _, text := range textQueue[comp.pos:] {
			if err := comp.session.SendData([]byte(text)); err != nil {
				log.Error(service.Current(ctx), err)
			}
		}

		comp.pos = len(textQueue)
	})
}

func (comp *HelloWorldComp) Shut() {
	runtime.Current(comp).Terminate()
}

func (comp *HelloWorldComp) onRecvData(session gate.ISession, data []byte) error {
	textMutex.Lock()
	defer textMutex.Unlock()

	text := fmt.Sprintf("[%s]:%s", comp.session.GetId(), string(data))
	textQueue = append(textQueue, text)

	log.Infof(service.Current(comp), text)
	return nil
}
