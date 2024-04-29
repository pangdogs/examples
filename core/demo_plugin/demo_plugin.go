package main

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/types"
	"git.golaxy.org/framework/plugins/log"
)

// demoPlugin 定义demo插件
var demoPlugin = define.ServicePlugin(func(...any) IDemoPlugin {
	return &DemoPlugin{}
})

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
	log.Infof(ctx, "init service plugin <%s>:[%s]", demoPlugin.Name, types.AnyFullName(*d))
	d.ctx = ctx
}

// ShutSP 关闭服务插件
func (d *DemoPlugin) ShutSP(ctx service.Context) {
	log.Infof(ctx, "shut service plugin <%s>:[%s]", demoPlugin.Name, types.AnyFullName(*d))
}

func (d *DemoPlugin) HelloWorld() {
	log.Infof(d.ctx, "plugin %q say hello world", demoPlugin.Name)
}
