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
	"time"

	"git.golaxy.org/core"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/iface"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"go.uber.org/zap"
)

/*
 * 基于core层提供的支持，演示一个简单的网关服务。
 */
func main() {
	// 创建服务并开始运行
	<-core.NewService(service.NewContext(
		service.With.Name("helloworld"),
		service.With.InstanceFace(iface.NewFaceT[service.Context](&HelloWorldService{})),
		service.With.RunningEventCB(func(svcCtx service.Context, runningEvent service.RunningEvent, _ ...any) {
			switch runningEvent {
			case service.RunningEvent_Birth:
				// 声明实体原型
				core.BuildEntityPT(svcCtx, "helloworld").
					AddComponent(HelloWorldComp{}).
					Declare()

				// 安装日志插件
				Log.Install(svcCtx)

				// 安装网关插件
				Gate.Install(svcCtx,
					GateWith.IOTimeout(3*time.Second),
					GateWith.IOBufferCap(1024*1024*5),
					GateWith.AgreeClientEncryptionProposal(true),
					GateWith.AgreeClientCompressionProposal(true),
					GateWith.CompressionThreshold(128),
					GateWith.SessionInactiveTimeout(10*time.Second),
				)

			case service.RunningEvent_Started:
				s := svcCtx.(*HelloWorldService)

				_, err := Gate.Require(svcCtx).Watch(svcCtx, generic.CastDelegateVoid1(s.handleSessionEstablished))
				if err != nil {
					log.L(svcCtx).Panic("watch session established failed", zap.Error(err))
				}

				log.L(svcCtx).Info("service started")

			case service.RunningEvent_Terminated:
				log.L(svcCtx).Info("service terminated")
			}
		}),
	)).Run().Done()
}

type HelloWorldService struct {
	service.ContextBehavior
}

func (s *HelloWorldService) handleSessionEstablished(session gate.ISession) {
	// 创建运行时并开始运行
	rt := core.NewRuntime(
		runtime.NewContext(s),
		core.With.Runtime.AutoRun(true),
	)

	// 在运行时中创建实体
	core.CallVoidAsync(rt, func(rtCtx runtime.Context, _ ...any) {
		entity, err := core.BuildEntity(rtCtx, "helloworld").
			SetPersistId(session.Id()).
			New()
		if err != nil {
			log.L(s).Panic("create entity failed", zap.Error(err))
		}
		log.L(s).Info("entity created", zap.String("entity_id", entity.Id().String()))

		go func() {
			<-session.Closed().Done()
			core.CallVoidAsync(entity, func(runtime.Context, ...any) { entity.Destroy() })
			<-entity.Terminated().Done()
			log.L(s).Info("entity destroyed", zap.String("entity_id", entity.Id().String()))
		}()
	})
}
