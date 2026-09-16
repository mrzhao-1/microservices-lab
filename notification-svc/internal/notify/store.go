package notify

import "sync"

// Notification 是一条通知记录。
type Notification struct {
	Event    string `json:"event"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Store 是通知的内存存储，v2 换成数据库。
type Store struct {
	mu    sync.Mutex
	items []Notification
}

func NewStore() *Store { return &Store{} }

func (s *Store) Add(n Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, n)
}

func (s *Store) List() []Notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Notification, len(s.items))
	copy(out, s.items)
	return out
}
