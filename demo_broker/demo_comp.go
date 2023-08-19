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
	sub, err := broker.Subscribe(service.Current(comp), context.Background(), "demo.>",
		broker.Option{}.EventHandler(func(e broker.Event) error {
			logger.Infof(service.Current(comp), "pattern:%s, topic:%s, receive: %s", e.Pattern(), e.Topic(), string(e.Message()))
			return nil
		}))
	if err != nil {
		logger.Panic(service.Current(comp), err)
	}
	comp.sub = sub

	logger.Infof(service.Current(comp), "max payload: %d", broker.MaxPayload(service.Current(comp)))

	golaxy.Await(comp, golaxy.AsyncTimeTick(service.Current(comp), time.Duration(rand.Int63n(5000))*time.Millisecond),
		func(ctx runtime.Context, ret runtime.Ret) {
			msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)
			err := broker.Publish(service.Current(comp), context.Background(), "demo.broker_test", []byte(msg))
			if err != nil {
				logger.Panic(service.Current(comp), err)
			}
			logger.Infof(service.Current(comp), "send: %s", msg)
			comp.sequence++
		})
}

// Shut 组件结束
func (comp *DemoComp) Shut() {
	comp.sub.Unsubscribe()
}
