package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"msg/api/internal/config"
	"msg/trans/transclient"
)

type ServiceContext struct {
	Config   config.Config
	TransRpc transclient.Trans
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		TransRpc: transclient.NewTrans(zrpc.MustNewClient(c.TransRpc)),
	}
}
