package main

import (
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/distributed"
	"kit.golaxy.org/plugins/log"
	"kit.golaxy.org/plugins/rpc"
	"kit.golaxy.org/plugins/rpc/callpath"
	"math/rand"
	"time"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

func (comp *DemoComp) Start() {
	rt := runtime.Current(comp)
	serv := service.Current(rt)

	golaxy.Await(rt, golaxy.TimeTick(rt, time.Second)).
		Pipe(rt, func(ctx runtime.Context, _ runtime.Ret, _ ...any) {
			var entityId string

			entities.AutoRLock(func(es *[]string) {
				if len(*es) <= 0 {
					return
				}
				entityId = (*es)[rand.Intn(len(*es))]
			})

			if entityId == comp.GetEntity().GetId().String() {
				return
			}

			dst := distributed.Using(serv).GetAddress().LocalAddr
			cp := callpath.CallPath{
				Category:  callpath.Entity,
				EntityId:  entityId,
				Component: "DemoComp",
				Method:    "HelloWorld",
			}

			a := rand.Int31()

			golaxy.Await(rt, rpc.RPC(serv, dst, cp.String(), a)).
				Any(rt, func(ctx runtime.Context, ret runtime.Ret, _ ...any) {
					rv, err := rpc.Result(ret)
					if err != nil {
						log.Errorf(serv, "3rd => result: %v", err)
						return
					}
					log.Infof(serv, "3rd => result: %d", rv)
				})

			log.Infof(service.Current(comp), "1st => call: %d", a)
		})
}

func (comp *DemoComp) HelloWorld(a int) int32 {
	n := rand.Int31()
	log.Infof(service.Current(comp), "2nd => accept: %d, return: %d", a, n)
	return n
}
