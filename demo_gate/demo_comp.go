package main

import (
	"fmt"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/gtp_gate"
	"kit.golaxy.org/plugins/logger"
	"sync"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[IDemoComp, DemoComp]("Demo组件")

// IDemoComp Demo组件接口
type IDemoComp interface {
	GetSession() gtp_gate.Session
}

// IDemoCompConstructor Demo组件构造函数
type IDemoCompConstructor interface {
	Constructor(session gtp_gate.Session)
}

var (
	textQueue []string
	textMutex sync.RWMutex
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	session gtp_gate.Session
	pos     int
}

func (comp *DemoComp) Start() {
	textMutex.RLock()
	defer textMutex.RUnlock()

	comp.pos = len(textQueue)

	golaxy.Await(runtime.Current(comp), golaxy.AsyncTimeTick(runtime.Current(comp), time.Second), func(ctx runtime.Context, ret runtime.Ret) {
		textMutex.RLock()
		defer textMutex.RUnlock()

		for _, text := range textQueue[comp.pos:] {
			if err := comp.session.SendData([]byte(text)); err != nil {
				logger.Error(service.Current(ctx), err)
			}
		}
		comp.pos = len(textQueue)
	})
}

func (comp *DemoComp) Shut() {
	runtime.Current(comp).GetCancelFunc()()
}

func (comp *DemoComp) Constructor(session gtp_gate.Session) {
	comp.session = session

	err := session.Options(gtp_gate.Option{}.SessionOption.RecvDataHandlers(comp.RecvDataHandler))
	if err != nil {
		logger.Panic(service.Current(comp), err)
	}
}

func (comp *DemoComp) RecvDataHandler(data []byte) error {
	textMutex.Lock()
	defer textMutex.Unlock()
	text := fmt.Sprintf("[%s]:%s", comp.session.GetId(), string(data))
	textQueue = append(textQueue, text)
	logger.Infof(service.Current(comp), text)
	return nil
}

func (comp *DemoComp) GetSession() gtp_gate.Session {
	return comp.session
}
