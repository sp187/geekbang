package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sp187/geekbang/final/internal/biz"
	fw "github.com/sp187/geekbang/final/internal/framework"
)

// UserPGWithRedisRepo pg+redis。
type UserPGWithRedisRepo struct {
	*UserPGRepo
	redis fw.AppCache
}

// NewUserPGWithRedisRepo 新建一个redis+pg的存储
func NewUserPGWithRedisRepo(repo *UserPGRepo, redis fw.AppCache) *UserPGWithRedisRepo {
	return &UserPGWithRedisRepo{
		UserPGRepo: repo,
		redis:      redis,
	}
}

func (u *UserPGWithRedisRepo) GetById(ctx context.Context, s uint) (biz.User, error) {
	key := fmt.Sprintf("%d", s)
	if b, err := u.redis.Get(ctx, key); err == nil {
		user := biz.User{}
		err = json.Unmarshal([]byte(b), &user)
		return user, err
	}
	user, err := u.UserPGRepo.GetById(ctx, s)
	if err != nil {
		err = u.redis.Set(ctx, key, user, 30*time.Minute)
		if err != nil {
			fw.GetLogger().Warn("set cache %s, %v fail", s, user)
		}
	}
	return user, nil
}

func (u *UserPGWithRedisRepo) Add(ctx context.Context, user biz.User) error {
	err := u.UserPGRepo.Add(ctx, user)
	if err != nil {
		return err
	}
	err = u.redis.Set(ctx, fmt.Sprintf("%d", user.Id), user, 30*time.Minute)
	if err != nil {
		fw.GetLogger().Warn("set cache %s, %v fail", user.Id, user)
	}
	return nil
}
