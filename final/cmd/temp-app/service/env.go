package service

import (
	"github.com/sp187/geekbang/final/internal/biz"
)

// 用户服务对象
var userService *biz.UserService

// InitAll 初始化必要的变量
func InitAll() error {
	var err error
	userService, err = NewService()
	if err != nil {
		panic("init user service fail")
	}
	return nil
}
