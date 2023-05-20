package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util"
	"kit.golaxy.org/plugins/logger"
)

// defineDemoPlugin 定义demo插件
var defineDemoPlugin = define.DefineServicePlugin[IDemoPlugin, any](func(options ...any) IDemoPlugin {
	return &DemoPlugin{
		options: options,
	}
})

// IDemoPlugin demo插件接口
type IDemoPlugin interface {
	HelloWorld()
}

// DemoPlugin demo插件实现
type DemoPlugin struct {
	options []any
	ctx     service.Context
}

// InitSP 初始化服务插件
func (d *DemoPlugin) InitSP(ctx service.Context) {
	logger.Infof(ctx, "init service plugin %q with %q", defineDemoPlugin.Name, util.TypeOfAnyFullName(*d))
	d.ctx = ctx
}

// ShutSP 关闭服务插件
func (d *DemoPlugin) ShutSP(ctx service.Context) {
	logger.Infof(ctx, "shut service plugin %q", defineDemoPlugin.Name)
}

func (d *DemoPlugin) HelloWorld() {
	logger.Infof(d.ctx, "%q say hello world", defineDemoPlugin.Name)
}
