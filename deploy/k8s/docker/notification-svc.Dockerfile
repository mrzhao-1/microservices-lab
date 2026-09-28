# notification-svc（go-micro）镜像。build context = 项目根目录。

FROM golang:1.27 AS build
WORKDIR /src
COPY notification-svc/ ./notification-svc/
WORKDIR /src/notification-svc
RUN CGO_ENABLED=0 go build -o /out/notification-svc ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/notification-svc ./notification-svc
EXPOSE 8081
ENTRYPOINT ["./notification-svc"]
