# Kubernetes 跑通方案

在容器/K8s 环境（如 `nzuguem/kubernetes-sandbox-environment`，内含 Docker + kind 集群）里跑通 microservices-lab。

## 与 Windows 版的区别

| 项 | Windows | K8s |
|---|---|---|
| 服务间地址 | localhost（代码默认值） | Service DNS（环境变量注入） |
| 前端托管 | 项目内 nginx.exe | 官方 nginx 镜像 |
| 对外暴露 | `http://localhost:80` | NodePort `30080` |

## 目录

- `docker/`：4 个 Dockerfile（gateway / user-svc / notification-svc / nginx）
- `nginx.conf`：K8s 版 nginx 配置（`proxy_pass http://gateway:8888`）
- `manifests/`：4 个 YAML（Deployment + Service）

## 步骤（在项目根目录执行）

1. 构建 4 个镜像：

```bash
docker build -f deploy/k8s/docker/gateway.Dockerfile -t microservices-lab/gateway:latest .
docker build -f deploy/k8s/docker/user-svc.Dockerfile -t microservices-lab/user-svc:latest .
docker build -f deploy/k8s/docker/notification-svc.Dockerfile -t microservices-lab/notification-svc:latest .
docker build -f deploy/k8s/docker/nginx.Dockerfile -t microservices-lab/nginx:latest .
```

2. 把镜像加载进 kind 集群（kind 的节点是容器，必须先 load）：

```bash
kind load docker-image \
  microservices-lab/gateway:latest \
  microservices-lab/user-svc:latest \
  microservices-lab/notification-svc:latest \
  microservices-lab/nginx:latest
```

3. 部署：

```bash
kubectl apply -f deploy/k8s/manifests/
```

4. 确认 Pod 全部 Running：

```bash
kubectl get pods
```

5. 访问：浏览器打开 `http://<节点IP>:30080/`（节点 IP 用 `kubectl get nodes -o wide` 查）。

## 验证链路

注册走通：gateway → user-svc 落库 → notification-svc 收到事件。

```bash
curl -X POST http://<节点IP>:30080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"zhaohu","email":"zhaohu@example.com"}'

# 看 notification-svc 是否收到事件
kubectl logs -l app=notification-svc
```

- 注册成功返回用户信息；同名再次注册返回 409。
- notification-svc 日志里应出现「收到事件 user.registered」。

## 清理

```bash
kubectl delete -f deploy/k8s/manifests/
```
