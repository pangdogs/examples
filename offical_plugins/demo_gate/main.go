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
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/meta"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/framework/plugins/gate"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/log/console_log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing endpoints")
	}

	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
	entityLib.Declare("demo", pt.Attribute{}, DemoComp{})

	// 创建插件包，安装插件
	pluginBundle := plugin.NewPluginBundle()
	console_log.Install(pluginBundle, console_log.With.Level(log.DebugLevel))

	// 安装网关插件
	gate.Install(pluginBundle,
		gate.With.TCPAddress(os.Args[1]),
		gate.With.IOTimeout(3*time.Second),
		gate.With.IOBufferCap(1024*1024*5),
		gate.With.AgreeClientEncryptionProposal(true),
		gate.With.AgreeClientCompressionProposal(true),
		gate.With.CompressedSize(128),
		gate.With.SessionInactiveTimeout(time.Hour),
		gate.With.SessionStateChangedHandler(generic.MakeDelegateAction3(handleSessionStateChanged)),
	)

	// 创建服务上下文与服务，并开始运行
	<-core.NewService(service.NewContext(
		service.With.EntityLib(entityLib),
		service.With.PluginBundle(pluginBundle),
		service.With.Name("demo_server"),
		service.With.RunningHandler(generic.MakeDelegateAction2(func(ctx service.Context, state service.RunningState) {
			if state != service.RunningState_Started {
				return
			}

			// 监听退出信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			go func() {
				<-sigChan
				ctx.Terminate()
			}()
		})),
	)).Run()
}

func handleSessionStateChanged(session gate.ISession, old, new gate.SessionState) {
	switch new {
	case gate.SessionState_Confirmed:
		// 创建运行时上下文与运行时，并开始运行
		rt := core.NewRuntime(runtime.NewContext(session.GetContext()),
			core.With.Runtime.AutoRun(true),
		)

		// 在运行时线程环境中，创建实体
		core.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
			entity, err := core.CreateEntity(ctx, "demo").
				Scope(ec.Scope_Global).
				PersistId(session.GetId()).
				Meta(meta.Make().Add("session", session).Get()).
				Spawn()
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}
			log.Infof(service.Current(ctx), "create entity %q finish", entity)
		}).Wait(session.GetContext())

	case gate.SessionState_Death:
		session.GetContext().CallVoid(uid.Id(session.GetId()), func(entity ec.Entity, _ ...any) {
			entity.DestroySelf()
		})
	}
}
