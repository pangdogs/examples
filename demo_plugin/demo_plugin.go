package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util/types"
	"kit.golaxy.org/plugins/log"
)

// demoPlugin 定义demo插件
var demoPlugin = define.DefineServicePlugin(func(...any) IDemoPlugin {
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
	log.Infof(ctx, "init service plugin %q with %q", demoPlugin.Name, types.AnyFullName(*d))
	d.ctx = ctx
}

// ShutSP 关闭服务插件
func (d *DemoPlugin) ShutSP(ctx service.Context) {
	log.Infof(ctx, "shut service plugin %q", demoPlugin.Name)
}

func (d *DemoPlugin) HelloWorld() {
	log.Infof(d.ctx, "plugin %q say hello world", demoPlugin.Name)
}
