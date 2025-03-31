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
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"log"
)

// 定义HelloWorld插件
var (
	self      = define.ServiceAddIn(func(...any) IHelloWorldAddIn { return &HelloWorldAddIn{} })
	Name      = self.Name
	Using     = self.Using
	Install   = self.Install
	Uninstall = self.Uninstall
)

// IHelloWorldAddIn HelloWorld插件接口
type IHelloWorldAddIn interface {
	HelloWorld()
}

// HelloWorldAddIn HelloWorld插件实现
type HelloWorldAddIn struct {
	svcCtx service.Context
}

// Init 初始化插件
func (h *HelloWorldAddIn) Init(svcCtx service.Context, _ runtime.Context) {
	log.Printf("init add-in %s", Name)
	h.svcCtx = svcCtx
}

// Shut 关闭插件
func (h *HelloWorldAddIn) Shut(ctx service.Context, _ runtime.Context) {
	log.Printf("shut add-in %s", Name)
}

// HelloWorld 打印HelloWorld
func (h *HelloWorldAddIn) HelloWorld() {
	log.Printf("add-in %s say hello world", Name)
}
