package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go-micro.dev/v5/web"

	"github.com/mrzhao-1/microservices-lab/notification-svc/internal/notify"
)

func main() {
	// go-micro 承载 HTTP 服务用 web.NewService（不是 micro.NewService，后者是 RPC）。
	// 本服务只做异步通知，是典型的「杂活」服务，适合用 go-micro 这类可插拔框架。
	svc := web.NewService(
		web.Name("notification-svc"),
		web.Address(":8081"),
	)

	store := notify.NewStore()

	svc.HandleFunc("/notify", handleNotify(store))
	svc.HandleFunc("/notifications", handleList(store))
	svc.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}

// handleNotify 收「用户已注册」事件，记录一条欢迎通知。
func handleNotify(store *notify.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var n notify.Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		store.Add(n)
		log.Printf("[notification] 收到事件 %s：给 %s(%s) 发欢迎通知", n.Event, n.Username, n.Email)
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// handleList 列出已记录的通知，方便自测。
func handleList(store *notify.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, store.List())
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
