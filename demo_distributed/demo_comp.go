package main

import (
	"encoding/json"
	"github.com/segmentio/ksuid"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/distributed"
	"kit.golaxy.org/plugins/gap/variant"
	"kit.golaxy.org/plugins/log"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

func (comp *DemoComp) Start() {
	golaxy.Await(runtime.Current(comp), golaxy.TimeTick(runtime.Current(comp), time.Second)).
		Pipe(runtime.Current(comp), func(ctx runtime.Context, ret runtime.Ret, _ ...any) {
			addr := distributed.Using(service.Current(ctx)).GetAddress()

			vmap, err := variant.MakeMap(map[string]int{
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
				ksuid.New().String(): rand.Int(),
			})
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			arr, err := variant.MakeArray([]int{rand.Int(), rand.Int(), rand.Int()})
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			msg := &MsgDemo{
				Int:    rand.Int(),
				Double: rand.Float64(),
				Str:    ksuid.New().String(),
				Map:    vmap,
				Array:  arr,
			}

			// 广播消息
			err = distributed.Using(service.Current(ctx)).SendMsg(addr.ServiceBroadcastAddr, msg)
			if err != nil {
				log.Panic(service.Current(ctx), err)
			}

			msgData, _ := json.Marshal(msg)
			log.Infof(service.Current(ctx), "send => topic:%q, msg:%s", addr.ServiceBroadcastAddr, msgData)
		})
}
