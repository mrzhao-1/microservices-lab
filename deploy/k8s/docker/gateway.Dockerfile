# gateway（go-zero）镜像。build context = 项目根目录。
# gateway 是独立 Go module，只 COPY 自己的目录。

FROM golang:1.27 AS build
WORKDIR /src
COPY gateway/ ./gateway/
WORKDIR /src/gateway
RUN CGO_ENABLED=0 go build -o /out/gateway .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/gateway ./gateway
COPY --from=build /src/gateway/etc/ ./etc/
EXPOSE 8888
ENTRYPOINT ["./gateway", "-f", "etc/gateway.yaml"]
