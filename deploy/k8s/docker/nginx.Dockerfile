# nginx 镜像：先构建前端（Vite），再用官方 nginx 镜像托管 + 反代。
# build context = 项目根目录。

FROM node:20-alpine AS febuild
WORKDIR /src
COPY frontend/ ./frontend/
WORKDIR /src/frontend
RUN npm ci && npm run build

FROM nginx:1.27-alpine
COPY deploy/k8s/nginx.conf /etc/nginx/nginx.conf
COPY --from=febuild /src/frontend/dist /usr/share/nginx/html
EXPOSE 80
