package main

import (
	"context"
	"fmt"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/broker"
	"kit.golaxy.org/plugins/logger"
	"math/rand"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[any, DemoComp]("Demo组件")

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	sub      broker.Subscriber
	sequence int
}

// Start 组件开始
func (comp *DemoComp) Start() {
	sub, err := broker.Subscribe(service.Get(comp), context.Background(), "demo.broker_test",
		broker.WithOption{}.EventHandler(func(e broker.Event) error {
			logger.Infof(service.Get(comp), "receive: %s", string(e.Message()))
			return nil
		}))
	if err != nil {
		logger.Panic(service.Get(comp), err)
	}
	comp.sub = sub

	logger.Infof(service.Get(comp), "max payload: %d", broker.MaxPayload(service.Get(comp)))

	golaxy.Await(comp, golaxy.AsyncTimeTick(service.Get(comp), time.Duration(rand.Int63n(5000))*time.Millisecond),
		func(ctx runtime.Context, ret runtime.Ret) {
			msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)
			err := broker.Publish(service.Get(comp), context.Background(), "demo.broker_test", []byte(msg))
			if err != nil {
				logger.Panic(service.Get(comp), err)
			}
			logger.Infof(service.Get(comp), "send: %s", msg)
			comp.sequence++
		})
}

// Shut 组件结束
func (comp *DemoComp) Shut() {
	comp.sub.Unsubscribe()
}
