package svc

import (
	"os"

	"github.com/mrzhao-1/microservices-lab/gateway/internal/config"
)

// ServiceContext 是贯穿整个网关生命周期的上下文，handler 从这里拿配置和下游地址。
type ServiceContext struct {
	Config config.Config
	// UserSvcAddr 是下游 user-svc（go-kratos）的 HTTP 地址。默认 localhost:8000，
	// 可用环境变量 USER_SVC_ADDR 覆盖（K8s 里指向 user-svc Service）。
	UserSvcAddr string
}

func NewServiceContext(c config.Config) *ServiceContext {
	addr := os.Getenv("USER_SVC_ADDR")
	if addr == "" {
		addr = "http://localhost:8000"
	}
	return &ServiceContext{
		Config:      c,
		UserSvcAddr: addr,
	}
}
