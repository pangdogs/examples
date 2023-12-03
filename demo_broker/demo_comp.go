package main

import (
	"context"
	"fmt"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/golaxy/util/generic"
	"kit.golaxy.org/plugins/broker"
	"kit.golaxy.org/plugins/log"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	sub      broker.Subscriber
	sequence int
}

// Start 组件开始
func (comp *DemoComp) Start() {
	log.Infof(service.Current(comp), "max payload: %d", broker.MaxPayload(service.Current(comp)))

	sub, err := broker.Subscribe(service.Current(comp), context.Background(), "demo.>",
		broker.Option{}.EventHandler(generic.CastDelegateFunc1(func(e broker.Event) error {
			log.Infof(service.Current(comp), "receive=> pattern:%q, topic:%q, msg:%q", e.Pattern(), e.Topic(), string(e.Message()))
			return nil
		})))
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
	comp.sub = sub

	golaxy.Await(runtime.Current(comp),
		golaxy.TimeTick(runtime.Current(comp), time.Duration(rand.Int63n(5000))*time.Millisecond),
	).Pipe(comp, func(ctx runtime.Context, _ runtime.Ret, _ ...any) {
		topic := "demo.broker_test"
		msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)

		if err := broker.Publish(service.Current(comp), context.Background(), topic, []byte(msg)); err != nil {
			log.Panic(service.Current(comp), err)
		}

		log.Infof(service.Current(comp), "send=> topic:%q, msg:%q", topic, msg)
		comp.sequence++
	})
}

// Shut 组件结束
func (comp *DemoComp) Shut() {
	comp.sub.Unsubscribe()
}
