package main

import (
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/log"
)

var (
	entityId = uid.New()
)

type HelloWorldService struct {
	framework.ServiceBehavior
}

func (s *HelloWorldService) Built(svc framework.IService) {
	s.BuildEntityPT("helloworld").
		AddComponent(HelloWorldComp{}).
		Declare()
}

func (s *HelloWorldService) Started(svc framework.IService) {
	entity, err := s.BuildEntityAsync("helloworld").
		SetPersistId(entityId).
		New()
	if err != nil {
		log.Panic(s, err)
	}

	go func() {
		<-entity.Terminated()
		<-s.Terminate()
	}()
}
