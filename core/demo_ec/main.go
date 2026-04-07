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
	"log"

	"git.golaxy.org/core"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
)

/*
 * 基于core层提供的支持，演示一个简单的EC系统实例，创建一个实体并运行，约10秒后结束。
 */
func main() {
	// 创建服务并开始运行
	<-core.NewService(service.NewContext(
		service.With.RunningEventCB(func(svcCtx service.Context, runningEvent service.RunningEvent, _ ...any) {
			switch runningEvent {
			case service.RunningEvent_Birth:
				// 声明实体原型
				core.BuildEntityPT(svcCtx, "helloworld").
					AddComponent(HelloWorldComp{}).
					Declare()

			case service.RunningEvent_Started:
				// 创建运行时并开始运行
				rt := core.NewRuntime(
					runtime.NewContext(svcCtx),
					core.With.Runtime.Frame(core.With.Frame.TotalFrames(10), core.With.Frame.TargetFPS(1)),
					core.With.Runtime.AutoRun(true),
				)

				// 在运行时中创建实体
				core.CallVoidAsync(rt, func(rtCtx runtime.Context, _ ...any) {
					entity, err := core.BuildEntity(rtCtx, "helloworld").New()
					if err != nil {
						log.Panic(err)
					}
					log.Printf("[%s] entity created", entity.Id())

					go func() {
						<-entity.Terminated().Done()
						log.Printf("[%s] entity destroyed", entity.Id())
						<-svcCtx.Terminate().Done()
					}()
				})

				log.Println("service started")

			case service.RunningEvent_Terminated:
				log.Println("service terminated")
			}
		}),
	)).Run().Done()
}
