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
	"errors"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/framework/plugins/dsync"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	mutex dsync.IDistMutex
}

// Update 组件更新
func (comp *DemoComp) Update() {
	if comp.mutex != nil {
		return
	}

	mutex := dsync.Using(service.Current(comp)).NewMutex("demo_dsync_counter")
	if err := mutex.Lock(context.Background()); err != nil {
		log.Errorf(service.Current(comp), "lock failed: %s", err)
		return
	}
	comp.mutex = mutex

	log.Info(service.Current(comp), "lock")

	content, err := os.ReadFile("demo_dsync_counter.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic(service.Current(comp), err)
		}
	}

	n, _ := strconv.Atoi(string(content))
	n++

	log.Infof(service.Current(comp), "counter: %d", n)

	err = os.WriteFile("demo_dsync_counter.txt", []byte(strconv.Itoa(n)), os.ModePerm)
	if err != nil {
		log.Panic(service.Current(comp), err)
	}

	core.Await(runtime.Current(comp),
		core.TimeAfter(context.Background(), time.Duration(rand.Int63n(1000))*time.Millisecond),
	).Any(func(ctx runtime.Context, _ async.Ret, _ ...any) {
		if comp.mutex == nil {
			return
		}
		comp.mutex.Unlock(context.Background())
		comp.mutex = nil

		log.Info(service.Current(comp), "unlock")
	})
}
