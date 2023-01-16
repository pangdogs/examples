package main

import (
	"fmt"
	"github.com/golaxy-kit/components/helloworld"
	"github.com/golaxy-kit/golaxy/plugin"
	"github.com/golaxy-kit/golaxy/pt"
	"github.com/golaxy-kit/golaxy/runtime"
	"github.com/golaxy-kit/golaxy/service"
	"github.com/golaxy-kit/golaxy/util"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("DistributedDemo", []string{
		util.TypeFullName[helloworld.HelloWorld](),
		util.TypeFullName[_DemoComp](),
	})

	// 创建插件库，安装插件
	pluginBundle := plugin.NewPluginBundle()
	registry_etcd.InstallTo(pluginBundle, registry_etcd.EtcdOption.Endpoints("127.0.0.1:3379"))

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.ContextOption.EntityLib(entityLib),
		service.ContextOption.PluginBundle(pluginBundle),
		service.ContextOption.StartedCallback(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			runtime := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.ContextOption.StoppedCallback(func(runtime.Context) {
						serviceCtx.GetCancelFunc()()
					}),
				),
				golaxy.RuntimeOption.Frame(runtime.NewFrame(30, 100, false)),
				golaxy.RuntimeOption.EnableAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			runtime.GetRuntimeCtx().SafeCallNoRetNoWait(func() {
				entity, err := golaxy.EntityCreator().
					RuntimeCtx(runtime.GetRuntimeCtx()).
					Prototype("DistributedDemo").
					Accessibility(golaxy.TryGlobal).
					TrySpawn()
				if err != nil {
					panic(err)
				}

				fmt.Printf("create entity[%s:%d:%d] finish\n", entity.GetPrototype(), entity.GetID(), entity.GetSerialNo())
			})
		}),
	)).Run()
}
