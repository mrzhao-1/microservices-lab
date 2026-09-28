# user-svc（go-kratos）镜像。build context = 项目根目录。

FROM golang:1.27 AS build
WORKDIR /src
COPY user-svc/ ./user-svc/
WORKDIR /src/user-svc
RUN CGO_ENABLED=0 go build -o /out/user-svc ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/user-svc ./user-svc
EXPOSE 8000
ENTRYPOINT ["./user-svc"]
