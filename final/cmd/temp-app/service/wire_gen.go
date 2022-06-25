// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"fmt"
	"github.com/sp187/geekbang/final/internal/biz"
	"github.com/sp187/geekbang/final/internal/data"
	"github.com/sp187/geekbang/final/internal/framework"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func NewService() (*biz.UserService, error) {
	config := NewDefaultConfig()
	db := NewDB(config)
	userPGRepo := data.NewUserPGRepo(db)
	appCache := NewCache(config)
	userRepo := NewRepo(userPGRepo, appCache)
	bizUserService, err := biz.NewUserService(userRepo)
	if err != nil {
		return nil, err
	}
	return bizUserService, nil
}

// wire.go:

func NewDefaultConfig() *fw.Config {
	return fw.LoadConfig("/mnt/sdb1/go/src/gitlab.bj.sensetime.com/sense-remote/project/web-template//configs/config.yaml")
}

func NewDB(cfg *fw.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	return fw.NewPGOrm(dsn)
}

func NewCache(cfg *fw.Config) fw.AppCache {
	address := fmt.Sprintf("%s:%s", cfg.CacheHost, cfg.CachePort)
	return fw.NewRedisClient(address, cfg.CachePassword)
}

func NewRepo(repo *data.UserPGRepo, cache fw.AppCache) biz.UserRepo {
	return data.NewUserPGWithRedisRepo(repo, cache)
}
