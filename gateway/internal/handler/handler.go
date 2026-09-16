package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/mrzhao-1/microservices-lab/gateway/internal/svc"
)

// RegisterHandlers 挂路由。go-zero 的 rest.Server.AddRoutes 接受一组 rest.Route。
func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/ping", Handler: pingHandler()},
		{Method: http.MethodPost, Path: "/api/user/register", Handler: registerHandler(ctx)},
	})
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJson(w, map[string]string{"msg": "pong"})
	}
}

// RegisterReq 是网关这一层收的请求体。网关只做透传 + 校验，业务结构归 user-svc 管。
type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func registerHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]any{"error": "请求体不是合法 JSON"})
			return
		}
		if req.Username == "" || req.Email == "" {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]any{"error": "username 和 email 不能为空"})
			return
		}

		resp, status, err := forwardRegister(ctx, req)
		if err != nil {
			httpx.WriteJson(w, status, map[string]any{"error": err.Error()})
			return
		}
		httpx.OkJson(w, resp)
	}
}

// forwardRegister 把请求原样转发给下游 user-svc，并透传它的状态码和响应体，
// 不让下游的业务错误（比如 409 重名）被误报成 502。
func forwardRegister(ctx *svc.ServiceContext, req RegisterReq) (map[string]any, int, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	resp, err := http.Post(ctx.UserSvcAddr+"/api/v1/users", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, http.StatusBadGateway, errors.New("调用 user-svc 失败: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New(extractError(data))
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, http.StatusBadGateway, err
	}
	return out, resp.StatusCode, nil
}

func extractError(data []byte) string {
	var e map[string]any
	if json.Unmarshal(data, &e) == nil {
		if m, ok := e["error"].(string); ok {
			return m
		}
	}
	return string(data)
}
