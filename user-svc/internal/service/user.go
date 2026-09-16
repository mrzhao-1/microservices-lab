package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mrzhao-1/microservices-lab/user-svc/internal/biz"
)

// notification-svc 的地址。v2 换成 NATS 消息总线 + go-micro broker 订阅。
const notifyAddr = "http://localhost:8081/notify"

// UserService 是 service 层（对应 DDD 的 application 层），负责 DTO 转换，不做复杂业务。
type UserService struct {
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

type registerReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register 处理 POST /api/v1/users。
func (s *UserService) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "请求体不是合法 JSON"})
		return
	}

	u, err := s.uc.Register(req.Username, req.Email)
	if err != nil {
		switch err {
		case biz.ErrUserExists:
			writeJSON(w, http.StatusConflict, map[string]any{"error": err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
		return
	}

	// 注册成功后异步通知（fire-and-forget，不阻塞本次响应）。
	s.notify(u)

	writeJSON(w, http.StatusOK, map[string]any{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}

// notify 把「用户已注册」事件发给 notification-svc。
func (s *UserService) notify(u *biz.User) {
	go func() {
		body, _ := json.Marshal(map[string]any{
			"event":    "user.registered",
			"username": u.Username,
			"email":    u.Email,
		})
		_, _ = http.Post(notifyAddr, "application/json", bytes.NewReader(body))
	}()
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
