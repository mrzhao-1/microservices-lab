package data

import (
	"sync"

	"github.com/mrzhao-1/microservices-lab/user-svc/internal/biz"
)

// userRepo 是 biz.UserRepo 的内存实现。v2 换成 gorm + MySQL 时只改这个文件。
type userRepo struct {
	mu     sync.Mutex
	users  map[string]*biz.User
	nextID int64
}

func NewUserRepo() biz.UserRepo {
	return &userRepo{users: make(map[string]*biz.User)}
}

func (r *userRepo) Create(u *biz.User) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	u.ID = r.nextID
	r.users[u.Username] = u
	return u, nil
}

func (r *userRepo) GetByUsername(username string) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[username]
	if !ok {
		return nil, biz.ErrNotFound
	}
	return u, nil
}
