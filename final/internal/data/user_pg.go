package data

import (
	"context"
	xerrors "github.com/pkg/errors"
	"github.com/sp187/geekbang/final/internal/biz"
	"gorm.io/gorm"
	"time"
)

// UserPGRepo 使用pg数据库实现biz.UserRepo
type UserPGRepo struct {
	db *gorm.DB //
}

// NewUserPGRepo 构建一个UserPGRepo
func NewUserPGRepo(db *gorm.DB) *UserPGRepo {
	db.AutoMigrate(&pgUser{})
	return &UserPGRepo{db}
}

func (u *UserPGRepo) GetById(ctx context.Context, id uint) (biz.User, error) {
	user := pgUser{
		ID: id,
	}
	err := u.db.WithContext(ctx).First(&user).Error
	return user.ToDO(), err
}

func (u *UserPGRepo) Add(ctx context.Context, user biz.User) error {
	dbUser := ToPG(user)
	result := u.db.WithContext(ctx).Create(&dbUser)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return xerrors.New("create user fail")
	}
	return nil
}

// pgUser 定义用户信息在数据库的表结构
type pgUser struct {
	ID        uint   `gorm:"type:serial"`
	Name      string `gorm:"column:username"`
	Gender    int
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (data pgUser) TableName() string {
	return "user"
}

// ToDO 转化为业务领域对象
func (u *pgUser) ToDO() biz.User {
	gender := "male"
	if u.Gender == 2 {
		gender = "female"
	}
	return biz.User{
		Id:     int(u.ID),
		Name:   u.Name,
		Age:    u.Age,
		Gender: gender,
	}
}

// ToPG 业务领域对象转化为存储对象
func ToPG(user biz.User) pgUser {
	gender := 1
	if user.Gender == "female" {
		gender = 2
	}
	return pgUser{
		Name:   user.Name,
		Gender: gender,
		Age:    user.Age,
	}
}
