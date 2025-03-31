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
	"context"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework/addins/dsync"
	"git.golaxy.org/framework/addins/log"
	"math/rand"
	"time"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	core.Await(runtime.Current(comp),
		core.TimeTickAsync(context.Background(), time.Duration(rand.Int63n(1000))*time.Millisecond),
	).Foreach(func(ctx runtime.Context, _ async.Ret, _ ...any) {
		mutex := dsync.Using(service.Current(comp)).NewMutex("helloworld", dsync.With.Expiry(10*time.Second), dsync.With.TimeoutFactor(0.5))
		if err := mutex.Lock(context.Background()); err != nil {
			log.Errorf(service.Current(comp), "[%s] lock failed: %s", comp.GetId(), err)
			return
		}
		defer mutex.Unlock(context.Background())

		log.Infof(service.Current(comp), "[%s] locked", comp.GetId())

		time.Sleep(time.Duration(rand.Int63n(200)) * time.Millisecond)

		log.Infof(service.Current(comp), "[%s] unlocked", comp.GetId())
	})
}
