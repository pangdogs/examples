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
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/log/zap_log"
	"time"
)

/*
 * 基于core层提供的支持，演示一个简单的网关服务，约10秒后结束。
 */
func main() {
	// 创建服务并开始运行
	<-core.NewService(service.NewContext(
		service.With.Name("helloworld"),
		service.With.RunningStatusChangedCB(func(svcCtx service.Context, state service.RunningStatus, _ ...any) {
			switch state {
			case service.RunningStatus_Starting:
				// 声明实体原型
				core.BuildEntityPT(svcCtx, "helloworld").
					AddComponent(HelloWorldComp{}).
					Declare()

				// 安装日志插件
				zap_log.Install(svcCtx)

				// 安装网关插件
				gate.Install(svcCtx,
					gate.With.IOTimeout(3*time.Second),
					gate.With.IOBufferCap(1024*1024*5),
					gate.With.AgreeClientEncryptionProposal(true),
					gate.With.AgreeClientCompressionProposal(true),
					gate.With.CompressedSize(128),
					gate.With.SessionInactiveTimeout(time.Hour),
					gate.With.SessionStateChangedHandler(generic.CastDelegateVoid3(onSessionStateChanged)),
				)

			case service.RunningStatus_Started:
				log.Info(svcCtx, "service started.")

			case service.RunningStatus_Terminated:
				log.Info(svcCtx, "service terminated.")
			}
		}),
	)).Run()
}

func onSessionStateChanged(session gate.ISession, curState, lastState gate.SessionState) {
	svcCtx := session.GetServiceContext()

	switch curState {
	case gate.SessionState_Confirmed:
		// 创建运行时并开始运行
		rt := core.NewRuntime(
			runtime.NewContext(svcCtx),
			core.With.Runtime.AutoRun(true),
		)

		// 在运行时中创建实体
		core.CallVoidAsync(rt, func(rtCtx runtime.Context, _ ...any) {
			entity, err := core.BuildEntity(rtCtx, "helloworld").
				SetPersistId(session.GetId()).
				SetMeta(map[string]any{"session": session}).
				New()
			if err != nil {
				log.Panic(svcCtx, err)
			}
			log.Infof(svcCtx, "[%s] entity created.", entity.GetId())

			go func() {
				<-entity.Terminated()
				log.Infof(svcCtx, "[%s] entity destroyed.", entity.GetId())
			}()
		}).Wait(svcCtx)

	case gate.SessionState_Death:
		svcCtx.CallVoidAsync(session.GetId(), func(entity ec.Entity, _ ...any) {
			entity.DestroySelf()
		})
	}
}
