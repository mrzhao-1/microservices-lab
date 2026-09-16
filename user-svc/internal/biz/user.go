package biz

import "errors"

// User 是领域对象，跨层使用。
type User struct {
	ID       int64
	Username string
	Email    string
}

var (
	ErrUserExists = errors.New("用户名已存在")
	ErrNotFound   = errors.New("用户不存在")
)

// UserRepo 是业务层定义的接口，data 层去实现它。
// 依赖倒置：biz 只管「要能存能查」，不关心存到内存还是 MySQL。
type UserRepo interface {
	Create(u *User) (*User, error)
	GetByUsername(username string) (*User, error)
}

// UserUsecase 是业务用例层，放核心业务规则。
type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Register 是核心业务：先查重，再落库。
func (uc *UserUsecase) Register(username, email string) (*User, error) {
	if _, err := uc.repo.GetByUsername(username); err == nil {
		return nil, ErrUserExists
	}
	return uc.repo.Create(&User{Username: username, Email: email})
}

// Get 按用户名查询。
func (uc *UserUsecase) Get(username string) (*User, error) {
	return uc.repo.GetByUsername(username)
}
