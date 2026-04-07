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
	"sync"

	"git.golaxy.org/core"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/uid"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/rpc/rpcpcsr"
	"go.uber.org/zap"
)

var (
	entityId = uid.New()
)

/*
 * 基于core层提供的支持，演示实体间RPC，约10秒后结束。
 */
func main() {
	var wg sync.WaitGroup

	for range 5 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// 创建服务并开始运行
			<-core.NewService(service.NewContext(
				service.With.Name("helloworld"),
				service.With.RunningEventCB(func(svcCtx service.Context, runningEvent service.RunningEvent, _ ...any) {
					switch runningEvent {
					case service.RunningEvent_Birth:
						// 声明实体原型
						core.BuildEntityPT(svcCtx, "helloworld").
							AddComponent(HelloWorldComp{}).
							Declare()

						// 安装日志插件
						Log.Install(svcCtx)

						// 安装服务发现插件
						DiscoveryEtcd.Install(svcCtx)

						// 安装broker插件
						BrokerNats.Install(svcCtx)

						// 安装分布式同步插件
						DsyncEtcd.Install(svcCtx)

						// 安装分布式服务插件
						Dsvc.Install(svcCtx)

						// 安装RPC插件
						RPC.Install(svcCtx, RPCWith.Processors(rpcpcsr.NewServiceProcessor(nil, false)))

					case service.RunningEvent_Started:
						// 创建运行时并开始运行
						rt := core.NewRuntime(runtime.NewContext(svcCtx,
							runtime.With.RunningEventCB(func(rtCtx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
								switch runningEvent {
								case runtime.RunningEvent_Birth:
									// 安装日志插件
									Log.Install(rtCtx)

									// 安装RPC栈插件
									RPCStack.Install(rtCtx)
								}
							})),
							core.With.Runtime.Frame(core.With.Frame.TotalFrames(10), core.With.Frame.TargetFPS(1)),
							core.With.Runtime.AutoRun(true),
						)

						// 在运行时中创建实体
						core.CallVoidAsync(rt, func(rtCtx runtime.Context, _ ...any) {
							entity, err := core.BuildEntity(rtCtx, "helloworld").SetPersistId(entityId).New()
							if err != nil {
								log.L(svcCtx).Panic("create entity failed", zap.Error(err))
							}
							log.L(svcCtx).Info("entity created", zap.String("entity_id", entity.Id().String()))

							go func() {
								<-entity.Terminated().Done()
								log.L(svcCtx).Info("entity destroyed", zap.String("entity_id", entity.Id().String()))
								<-svcCtx.Terminate().Done()
							}()
						})

						log.L(svcCtx).Info("service started")

					case service.RunningEvent_Terminated:
						log.L(svcCtx).Info("service terminated")
					}
				}),
			)).Run().Done()
		}()
	}

	wg.Wait()
}
