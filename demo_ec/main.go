package main

import (
	"fmt"
	"github.com/golaxy-kit/components/helloworld"
	_ "github.com/golaxy-kit/components/setup"
	"github.com/golaxy-kit/golaxy"
	"github.com/golaxy-kit/golaxy/ec"
	"github.com/golaxy-kit/golaxy/pt"
	"github.com/golaxy-kit/golaxy/runtime"
	"github.com/golaxy-kit/golaxy/service"
)

func main() {
	// 创建实体库，注册实体原型
	entityLib := pt.NewEntityLib()
	entityLib.Register("demo", []string{
		helloworld.Comp.Name,
		DemoComp.Name,
	})

	// 创建服务上下文与服务，并开始运行
	<-golaxy.NewService(service.NewContext(
		service.WithContextOption{}.EntityLib(entityLib),
		service.WithContextOption{}.StartedCallback(func(serviceCtx service.Context) {
			// 创建运行时上下文与运行时，并开始运行
			rt := golaxy.NewRuntime(
				runtime.NewContext(serviceCtx,
					runtime.WithContextOption{}.StoppedCallback(func(runtime.Context) {
						serviceCtx.GetCancelFunc()()
					}),
				),
				golaxy.WithRuntimeOption{}.Frame(runtime.NewFrame(30, 300, false)),
				golaxy.WithRuntimeOption{}.EnableAutoRun(true),
			)

			// 在运行时线程环境中，创建实体
			rt.GetRuntimeCtx().SafeCallNoRetNoWait(func() {
				entity, err := golaxy.NewEntityCreator(rt.GetRuntimeCtx(),
					pt.WithEntityOption{}.Prototype("demo"),
					pt.WithEntityOption{}.Accessibility(ec.Global),
				).Spawn()
				if err != nil {
					panic(err)
				}

				fmt.Printf("create entity %s finish\n", entity)
			})
		}),
	)).Run()
}
