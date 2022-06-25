package data

import (
	"context"

	"github.com/sp187/geekbang/final/internal/biz"
)

// UserMemRepo 使用内存实现biz.UserRepo
type UserMemRepo struct {
	user map[uint]biz.User
}

func NewUserMemRepo() *UserMemRepo {
	return &UserMemRepo{user: make(map[uint]biz.User)}
}

// GetById ctx用于传递上下文信息。
func (u *UserMemRepo) GetById(ctx context.Context, s uint) (biz.User, error) {
	user, ok := u.user[s]
	if !ok {
		return user, biz.ErrUserNotFound
	}
	return user, nil
}

func (u *UserMemRepo) Add(ctx context.Context, user biz.User) error {
	u.user[uint(user.Id)] = user
	return nil
}
