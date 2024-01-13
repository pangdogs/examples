package main

import (
	"fmt"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/plugins/gtp_gate"
	"git.golaxy.org/plugins/log"
	"sync"
	"time"
)

var (
	textQueue []string
	textMutex sync.RWMutex
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	session gtp_gate.ISession
	pos     int
}

func (comp *DemoComp) Awake() {
	comp.session = comp.GetEntity().GetMeta().Get("session").(gtp_gate.ISession)
}

func (comp *DemoComp) Start() {
	textMutex.RLock()
	defer textMutex.RUnlock()

	comp.pos = len(textQueue)

	core.Await(runtime.Current(comp),
		core.TimeTick(runtime.Current(comp), time.Second),
	).Pipe(runtime.Current(comp), func(ctx runtime.Context, ret runtime.Ret, _ ...any) {
		textMutex.RLock()
		defer textMutex.RUnlock()

		for _, text := range textQueue[comp.pos:] {
			if err := comp.session.SendData([]byte(text)); err != nil {
				log.Error(service.Current(ctx), err)
			}
		}
		comp.pos = len(textQueue)
	})
}

func (comp *DemoComp) Shut() {
	runtime.Current(comp).GetCancelFunc()()
}

func (comp *DemoComp) Constructor(session gtp_gate.ISession) {
	comp.session = session

	err := session.Settings(gtp_gate.Option{}.Session.RecvDataHandler(generic.CastDelegateFunc1(comp.RecvDataHandler)))
	if err != nil {
		log.Panic(session.GetContext(), err)
	}
}

func (comp *DemoComp) RecvDataHandler(data []byte) error {
	textMutex.Lock()
	defer textMutex.Unlock()
	text := fmt.Sprintf("[%s]:%s", comp.session.GetId(), string(data))
	textQueue = append(textQueue, text)
	log.Infof(service.Current(comp), text)
	return nil
}

func (comp *DemoComp) GetSession() gtp_gate.ISession {
	return comp.session
}
