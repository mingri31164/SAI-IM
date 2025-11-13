package svc

import (
	"SAI-IM/apps/task/mq/internal/config"
)

type ServiceContext struct {
	config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config: c,
	}
	return svc
}
