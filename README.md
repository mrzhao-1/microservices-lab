# microservices-lab

一个练手用的微服务框架对比项目：前端经 nginx 路由，后端是三个**互相独立**的 Go 服务，分别用 go-zero、go-kratos、go-micro 三个框架实现，看它们各自擅长什么、边界怎么划。

## 分界（核心）

三个服务是三个独立的 Go 模块、三个独立的进程，各自单独编译、单独启动。它们的责任不重叠：

| 项目 | 框架 | 端口 | 责任 | 用到的招牌能力 |
|---|---|---|---|---|
| `gateway` | **go-zero** | 8888 | 唯一对外入口：鉴权、参数校验、转发 | rest 路由 + config + 中间件 |
| `user-svc` | **go-kratos** | 8000 | 核心业务：用户注册、查询 | service/biz/data 三层 + 依赖倒置 |
| `notification-svc` | **go-micro** | 8081 | 异步通知：注册成功后发欢迎通知 | Service 抽象 + client/registry 概念 |

## 为什么这么分

- **go-zero 放网关**：它的强项是 API 网关这一套，`.api` 定义 + 代码生成 + 限流熔断降载，本来就是给「边缘层」用的。
- **go-kratos 放核心业务**：它的招牌是 `service → biz → data` 的分层规范和依赖倒置，适合写有领域逻辑的业务服务。
- **go-micro 放异步通知**：它的招牌是可插拔的 Service/registry/broker 抽象，适合做事件、异步、多协议这种「杂活」。

## 数据流（v1，零外部依赖）

```
浏览器 ──HTTP──► nginx ──► gateway(go-zero) ──HTTP──► user-svc(kratos)
                                                         │ 注册成功
                                                         │ HTTP 异步通知(fire-and-forget)
                                                         ▼
                                       notification-svc(go-micro) ──记录通知──► 内存
```

## 目录结构

```
microservices-lab/
├── gateway/             go-zero 网关
├── user-svc/            go-kratos 核心业务服务
├── notification-svc/    go-micro 异步通知服务
├── frontend/            Vue3 + Vite 前端（一个注册页）
├── nginx/               nginx 1.26.3（已装进项目：托管前端 + 反代 /api）
├── deploy/              两种运行环境的跑通方案（并列放置）
│   ├── windows/         Windows 本地跑通说明
│   └── k8s/             Kubernetes 跑通（Dockerfile + 清单）
└── README.md
```

## 运行（两种环境）

跑通方案按环境并列放在 `deploy/` 下，整个项目是同一套代码，只是运行环境不同：

| 环境 | 说明 | 入口 |
|---|---|---|
| Windows 本地 | `go run` 起三个 Go 服务 + 项目内 nginx.exe | [deploy/windows/README.md](deploy/windows/README.md) |
| Kubernetes | Docker 镜像 + K8s 清单，kind 集群里跑 | [deploy/k8s/README.md](deploy/k8s/README.md) |

## 关于 nginx

`nginx/conf/nginx.conf` 是生效配置：`location /` 托管 `nginx/html` 里的前端产物（SPA 用 try_files），`location /api/` 反代到网关 8888。完整链路：

```
浏览器 ──► nginx(:80) ──► gateway(go-zero:8888) ──► user-svc(kratos:8000)
                                                       │
                                                       └──► notification-svc(go-micro:8081)
```

nginx 只在最外层做「静态托管 + 入口反代」，业务鉴权、参数校验、路由在 go-zero 网关，两层不冲突、上下各司其职。

## v2 升级点（这版故意不做，留作练习）

| 能力 | v1 现状 | v2 怎么升级 |
|---|---|---|
| 存储 | 内存 map | MySQL + gorm，或 SQLite |
| 异步 | HTTP 回调 | NATS 消息总线 + go-micro broker 订阅 |
| 服务发现 | 固定地址 | etcd（注意：三个框架的 registry 格式不互通，网关要自己适配） |
| 鉴权 | 网关里一个 stub | 真 JWT，user-svc 发 token |
| 协议 | HTTP 互通 | 网关 zRPC 直连 kratos 的 gRPC 端口 |

## 技术栈版本（钉死，避免飘）

- Go 1.27
- go-zero **v1.10.3**
- go-kratos **v2.9.2**（v2 最新版。注意：≤2.9.2 有公开漏洞，修复只在 v3.0.0，v2 线暂无补丁。漏洞在 HTTP 的 DefaultServeMux 兜底路径，本项目的 kratos 服务用显式 `HandleFunc` 路由，不碰那条路径；上线生产前再升 v3）
- go-micro **v5.30.0**（正确路径 `go-micro.dev/v5`，老仓库 `github.com/micro/go-micro` 已废弃。仅作学习）
- Vue 3 + Vite

> 状态：三个服务 `go build` 全通过，端到端冒烟测试已跑通（注册 → user-svc 落库 → notification-svc 收到事件；重名返回 409）。
