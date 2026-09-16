package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"

	"github.com/mrzhao-1/microservices-lab/user-svc/internal/biz"
	"github.com/mrzhao-1/microservices-lab/user-svc/internal/data"
	"github.com/mrzhao-1/microservices-lab/user-svc/internal/server"
	"github.com/mrzhao-1/microservices-lab/user-svc/internal/service"
)

func main() {
	// 手动依赖注入（不用 wire 生成器，看得更清楚）：
	// data 实现 biz 定义的接口 → biz 提供业务逻辑 → service 做 DTO 转换 → server 挂路由。
	repo := data.NewUserRepo()
	uc := biz.NewUserUsecase(repo)
	svc := service.NewUserService(uc)
	hs := server.NewHTTPServer(svc)

	app := kratos.New(
		kratos.Name("user-svc"),
		kratos.Version("v1.0.0"),
		kratos.Server(hs),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
