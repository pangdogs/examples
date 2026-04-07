package main

import (
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/framework"
	"go.uber.org/zap"
)

var (
	entityId = uid.New()
)

type HelloWorldService struct {
	framework.ServiceBehavior
}

func (s *HelloWorldService) OnBuilt(svc framework.IService) {
	s.BuildEntityPT("helloworld").
		AddComponent(HelloWorldComp{}).
		Declare()
}

func (s *HelloWorldService) OnStarted(svc framework.IService) {
	entity, err := s.BuildEntity("helloworld").
		SetPersistId(entityId).
		New()
	if err != nil {
		s.L().Panic("create entity failed", zap.Error(err))
	}

	go func() {
		<-entity.Terminated().Done()
		<-s.Terminate().Done()
	}()
}
