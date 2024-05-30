package main

import (
	"fmt"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/plugins/gate"
	"git.golaxy.org/framework/plugins/log"
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
	session gate.ISession
	pos     int
}

func (comp *DemoComp) Awake() {
	comp.session = comp.GetEntity().GetMeta().Get("session").(gate.ISession)
}

func (comp *DemoComp) Start() {
	textMutex.RLock()
	defer textMutex.RUnlock()

	comp.pos = len(textQueue)

	core.Await(runtime.Current(comp),
		core.TimeTick(runtime.Current(comp), time.Second),
	).Pipe(runtime.Current(comp), func(ctx runtime.Context, ret async.Ret, _ ...any) {
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
	runtime.Current(comp).Terminate()
}

func (comp *DemoComp) Constructor(session gate.ISession) {
	comp.session = session

	err := session.GetSettings().RecvDataHandler(generic.MakeDelegateFunc2(comp.RecvDataHandler)).Change()
	if err != nil {
		log.Panic(session.GetContext(), err)
	}
}

func (comp *DemoComp) RecvDataHandler(session gate.ISession, data []byte) error {
	textMutex.Lock()
	defer textMutex.Unlock()
	text := fmt.Sprintf("[%s]:%s", comp.session.GetId(), string(data))
	textQueue = append(textQueue, text)
	log.Infof(service.Current(comp), text)
	return nil
}

func (comp *DemoComp) GetSession() gate.ISession {
	return comp.session
}
