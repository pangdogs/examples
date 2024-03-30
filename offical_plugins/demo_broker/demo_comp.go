package main

import (
	"context"
	"fmt"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/framework/plugins/broker"
	"git.golaxy.org/framework/plugins/log"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	sub      broker.ISubscriber
	sequence int
}

// Start 组件开始
func (comp *DemoComp) Start() {
	log.Infof(service.Current(comp), "max payload: %d", broker.Using(service.Current(comp)).GetMaxPayload())

	sub, err := broker.Using(service.Current(comp)).Subscribe(context.Background(), "demo.>",
		broker.With.EventHandler(generic.CastDelegateFunc1(func(e broker.IEvent) error {
			log.Infof(service.Current(comp), "receive=> pattern:%q, topic:%q, msg:%q", e.Pattern(), e.Topic(), string(e.Message()))
			return nil
		})))
	if err != nil {
		log.Panic(service.Current(comp), err)
	}
	comp.sub = sub

	core.Await(runtime.Current(comp),
		core.TimeTick(runtime.Current(comp), time.Duration(rand.Int63n(5000))*time.Millisecond),
	).Pipe(nil, func(ctx runtime.Context, _ runtime.Ret, _ ...any) {
		topic := "demo.broker_test"
		msg := fmt.Sprintf("%s-%d", comp.GetId(), comp.sequence)

		if err := broker.Using(service.Current(comp)).Publish(context.Background(), topic, []byte(msg)); err != nil {
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
