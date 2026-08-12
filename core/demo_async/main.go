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
	"log"
	"time"

	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
)

type LoaderComp struct {
	ec.ComponentBehavior
	value string
}

// Start 将阻塞工作放到组件 Scope，再把结果送回 Runtime 更新组件状态。
func (comp *LoaderComp) Start() {
	// Spawn 中的函数运行在后台 goroutine，并在组件销毁时收到取消信号。
	load := core.Spawn(comp, func(ctx context.Context, _ ...any) async.Result {
		timer := time.NewTimer(250 * time.Millisecond)
		defer timer.Stop()

		select {
		case <-timer.C:
			return async.NewResult("loaded outside the Runtime", nil)
		case <-ctx.Done():
			return async.NewResult(nil, ctx.Err())
		}
	})

	// ContinueOnVoid 的回调由组件所属 Runtime 串行执行，可以安全修改组件状态。
	core.ContinueOnVoid(comp, load, func(_ runtime.Context, ret async.Result, _ ...any) {
		if ret.Error != nil {
			log.Printf("load failed: %v", ret.Error)
			comp.Entity().Destroy()
			return
		}

		comp.value = ret.Value.(string)
		log.Printf("result applied on Runtime: %s", comp.value)
		comp.Entity().Destroy()
	})
}

func main() {
	<-core.NewService(service.NewContext(
		service.With.RunningEventCB(func(svcCtx service.Context, runningEvent service.RunningEvent, _ ...any) {
			switch runningEvent {
			case service.RunningEvent_Birth:
				core.BuildEntityPT(svcCtx, "loader").
					AddComponent(LoaderComp{}).
					Declare()

			case service.RunningEvent_Started:
				rt := core.NewRuntime(
					runtime.NewContext(svcCtx),
					core.With.Runtime.AutoRun(true),
				)

				if err := core.Post(rt, func(rtCtx runtime.Context, _ ...any) {
					entity, err := core.BuildEntity(rtCtx, "loader").New()
					if err != nil {
						log.Panic(err)
					}

					go func() {
						<-entity.Terminated().Done()
						<-svcCtx.Terminate().Done()
					}()
				}); err != nil {
					log.Panicf("enqueue entity creation failed: %v", err)
				}
			}
		}),
	)).Run().Done()
}
