package main

import (
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/logger"
	"reflect"
)

// defineDemoPlugin 定义demo插件
var defineDemoPlugin = define.DefineServicePlugin[IDemoPlugin, any](func(options ...any) IDemoPlugin {
	return &_DemoPlugin{
		options: options,
	}
})

// IDemoPlugin demo插件接口
type IDemoPlugin interface {
	HelloWorld()
}

// _DemoPlugin demo插件实现
type _DemoPlugin struct {
	options []any
	ctx     service.Context
}

// InitSP 初始化服务插件
func (d *_DemoPlugin) InitSP(ctx service.Context) {
	logger.Infof(ctx, "init service plugin %q with %q", defineDemoPlugin.Name, reflect.TypeOf(d).Elem())
	d.ctx = ctx
}

// ShutSP 关闭服务插件
func (d *_DemoPlugin) ShutSP(ctx service.Context) {
	logger.Infof(ctx, "shut service plugin %q", defineDemoPlugin.Name)
}

func (d *_DemoPlugin) HelloWorld() {
	logger.Infof(d.ctx, "%q say hello world", defineDemoPlugin.Name)
}
