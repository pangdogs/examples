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
	"fmt"
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/types"
)

// demoPlugin 定义demo插件
var demoPlugin = define.ServicePlugin(func(...any) IDemoPlugin { return &DemoPlugin{} })

var (
	Using     = demoPlugin.Using
	Install   = demoPlugin.Install
	Uninstall = demoPlugin.Uninstall
)

// IDemoPlugin demo插件接口
type IDemoPlugin interface {
	HelloWorld()
}

// DemoPlugin demo插件实现
type DemoPlugin struct {
	ctx service.Context
}

// InitSP 初始化服务插件
func (d *DemoPlugin) InitSP(ctx service.Context) {
	fmt.Printf("init service plugin <%s>:[%s]\n", demoPlugin.Name, types.FullNameT[DemoPlugin]())
	d.ctx = ctx
}

// ShutSP 关闭服务插件
func (d *DemoPlugin) ShutSP(ctx service.Context) {
	fmt.Printf("shut service plugin <%s>:[%s]", demoPlugin.Name, types.FullNameT[DemoPlugin]())
}

// HelloWorld Hello World
func (d *DemoPlugin) HelloWorld() {
	fmt.Printf("plugin %q say hello world", demoPlugin.Name)
}
