# Windows 本地跑通方案

这是 microservices-lab 在 Windows 本机的运行方式：三个 Go 服务用 `go run` 直接起，前端用项目里自带的 nginx（1.26.3，Windows 版）托管。

## 前置

- Go 1.27
- nginx 已装进项目（`nginx/` 目录，含 `nginx.exe`），无需额外安装
- 前端产物已构建好（`nginx/html`），无需重新 build

## 启动（完整链路，含 nginx）

1. 起 user-svc：`cd user-svc && go run ./cmd/server`
2. 起 notification-svc：`cd notification-svc && go run ./cmd`
3. 起 gateway：`cd gateway && go run .`
4. 起 nginx：`cd nginx && nginx.exe`
5. 浏览器打开 `http://localhost/`，页面注册即可走通全链路

> 服务间地址用默认值（localhost），无需设环境变量。改过前端代码后，重新 `cd frontend && npm run build`，再把 `dist` 拷进 `nginx/html`。

## 直接调接口验证（跳过浏览器）

```
curl -X POST http://localhost/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"zhaohu","email":"zhaohu@example.com"}'
```

注册成功返回用户信息；同名再次注册返回 409。

## 开发模式（不想用 nginx 时）

前端用 Vite dev server 代理（`frontend/vite.config.js` 已配好，把 `/api` 转发到 8888）：

```
cd frontend && npm run dev
```
