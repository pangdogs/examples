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
	"git.golaxy.org/core"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"log"
)

/*
 * 基于core层提供的支持，演示一个简单的服务插件案例，约10秒后结束。
 */
func main() {
	// 创建服务并开始运行
	<-core.NewService(service.NewContext(
		service.With.RunningStatusChangedCB(func(svcCtx service.Context, state service.RunningStatus, _ ...any) {
			switch state {
			case service.RunningStatus_Starting:
				// 声明实体原型
				core.BuildEntityPT(svcCtx, "helloworld").
					AddComponent(HelloWorldComp{}).
					Declare()

				// 安装插件
				Install(svcCtx)

			case service.RunningStatus_Started:
				// 创建运行时并开始运行
				rt := core.NewRuntime(
					runtime.NewContext(svcCtx),
					core.With.Runtime.Frame(runtime.NewFrame(
						runtime.With.Frame.TotalFrames(10),
						runtime.With.Frame.TargetFPS(1),
					)),
					core.With.Runtime.AutoRun(true),
				)

				// 在运行时中创建实体
				core.CallVoidAsync(rt, func(rtCtx runtime.Context, _ ...any) {
					entity, err := core.BuildEntity(rtCtx, "helloworld").New()
					if err != nil {
						log.Panic(err)
					}
					log.Printf("[%s] entity created.", entity.GetId())

					go func() {
						<-entity.Terminated()
						log.Printf("[%s] entity destroyed.", entity.GetId())
						<-svcCtx.Terminate()
					}()
				})

				log.Println("service started.")

			case service.RunningStatus_Terminated:
				log.Println("service terminated.")
			}
		}),
	)).Run()
}
