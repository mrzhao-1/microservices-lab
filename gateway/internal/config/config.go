package config

import "github.com/zeromicro/go-zero/rest"

// Config 是网关的全部配置，内嵌 go-zero 的 RestConf（含 Host、Port 等）。
type Config struct {
	rest.RestConf
}
