package server

import (
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/mrzhao-1/microservices-lab/user-svc/internal/conf"
	"github.com/mrzhao-1/microservices-lab/user-svc/internal/service"
)

// NewHTTPServer 建一个 kratos 的 HTTP server，挂 service 提供的路由。
func NewHTTPServer(svc *service.UserService) *khttp.Server {
	srv := khttp.NewServer(khttp.Address(conf.HTTPAddr))

	srv.HandleFunc("/api/v1/users", svc.Register)
	srv.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	return srv
}
