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
	"math/rand"
	"time"

	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/service"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/dsync"
	"git.golaxy.org/framework/addins/log"
	"go.uber.org/zap"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	svcCtx := service.Current(comp)
	logger := log.L(svcCtx)
	interval := time.Duration(rand.Int63n(1000)+1) * time.Millisecond

	core.SpawnVoid(comp, func(ctx context.Context, _ ...any) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			mutex := Dsync.Require(svcCtx).NewMutex("helloworld", dsync.With.Expiry(10*time.Second), dsync.With.TimeoutFactor(0.5))
			if err := mutex.Lock(ctx); err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Error("lock failed", zap.Error(err))
				continue
			}

			logger.Info("locked")

			hold := time.NewTimer(time.Duration(rand.Int63n(200)) * time.Millisecond)
			select {
			case <-hold.C:
			case <-ctx.Done():
				if !hold.Stop() {
					select {
					case <-hold.C:
					default:
					}
				}
			}

			unlockCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err := mutex.Unlock(unlockCtx)
			cancel()
			if err != nil {
				logger.Error("unlock failed", zap.Error(err))
			} else {
				logger.Info("unlocked")
			}

			if ctx.Err() != nil {
				return
			}
		}
	})
}
