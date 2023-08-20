package main

import (
	"fmt"
	"kit.golaxy.org/golaxy"
	"kit.golaxy.org/golaxy/define"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/service"
	"kit.golaxy.org/plugins/gate"
	gtp_gate "kit.golaxy.org/plugins/gate/gtp"
	"kit.golaxy.org/plugins/logger"
	"sync"
	"time"
)

// defineDemoComp 定义Demo组件
var defineDemoComp = define.DefineComponent[IDemoComp, DemoComp]("Demo组件")

// IDemoComp Demo组件接口
type IDemoComp interface {
	GetSession() gate.Session
}

// IDemoCompConstructor Demo组件构造函数
type IDemoCompConstructor interface {
	Constructor(session gate.Session)
}

var (
	textQueue []string
	textMutex sync.RWMutex
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
	session gate.Session
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

func (comp *DemoComp) Constructor(session gate.Session) {
	setting, err := gtp_gate.GetSessionSetting(session)
	if err != nil {
		logger.Panic(service.Current(comp), err)
	}

	setting.RecvDataHandlers(func(session gate.Session, data []byte) error {
		textMutex.Lock()
		defer textMutex.Unlock()
		text := fmt.Sprintf("[%s]:%s", session.GetId(), string(data))
		textQueue = append(textQueue, text)
		logger.Infof(service.Current(comp), text)
		return nil
	})

	comp.session = session
}

func (comp *DemoComp) GetSession() gate.Session {
	return comp.session
}
