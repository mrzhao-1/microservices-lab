package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"github.com/mrzhao-1/microservices-lab/gateway/internal/config"
	"github.com/mrzhao-1/microservices-lab/gateway/internal/handler"
	"github.com/mrzhao-1/microservices-lab/gateway/internal/svc"
)

var configFile = flag.String("f", "etc/gateway.yaml", "配置文件路径")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	server.Start()
}
