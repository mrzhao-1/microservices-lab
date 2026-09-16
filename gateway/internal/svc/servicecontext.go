package svc

import "github.com/mrzhao-1/microservices-lab/gateway/internal/config"

// ServiceContext 是贯穿整个网关生命周期的上下文，handler 从这里拿配置和下游地址。
type ServiceContext struct {
	Config config.Config
	// UserSvcAddr 是下游 user-svc（go-kratos）的 HTTP 地址。v1 用固定地址，
	// v2 换成服务发现时只改这里。
	UserSvcAddr string
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		UserSvcAddr: "http://localhost:8000",
	}
}
