package main

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/framework/plugins/dserv"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/callpath"
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

	core.Await(rt, core.TimeTick(rt, 3*time.Second)).
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

			dst := dserv.Using(serv).GetAddress().LocalAddr

			addr := dserv.Using(serv).GetAddress()
			_ = addr

			cp1 := callpath.CallPath{
				Category:  callpath.Entity,
				EntityId:  entityId,
				Component: "DemoComp",
				Method:    "TestRPC",
			}

			a := rand.Uint32()

			// 异步
			{
				core.Await(rt, rpc.RPC(serv, dst, cp1.String(), a)).
					Any(rt, func(ctx runtime.Context, ret runtime.Ret, _ ...any) {
						rv, err := rpc.Result1[int32](ret)
						if err != nil {
							log.Errorf(serv, "3rd => result: %v", err)
							return
						}
						log.Infof(serv, "3rd => result: %d", rv)
					})
			}

			//// 同步
			//{
			//	rv, err := rpc.Result1[int32](<-rpc.RPC(serv, dst, cp1.String(), a))
			//	if err != nil {
			//		log.Errorf(serv, "3rd => result: %v", err)
			//	} else {
			//		log.Infof(serv, "3rd => result: %d", rv)
			//	}
			//}

			log.Infof(service.Current(comp), "1st => call: %d", a)

			cp2 := callpath.CallPath{
				Category:  callpath.Entity,
				EntityId:  entityId,
				Component: "DemoComp",
				Method:    "TestOneWayRPC",
			}

			err := rpc.OneWayRPC(serv, dst, cp2.String(), a)
			if err != nil {
				log.Errorf(serv, "oneway => call: %v", err)
				return
			}
		})
}

func (comp *DemoComp) TestRPC(a uint32) int32 {
	n := rand.Int31()
	log.Infof(service.Current(comp), "2nd => accept: %d, return: %d", a, n)
	return n
}

func (comp *DemoComp) TestOneWayRPC(a uint32) {
	log.Infof(service.Current(comp), "oneway => accept: %d", a)
}
