package biz

import (
	"context"

	errs "github.com/sp187/geekbang/final/internal/framework/error"
	"github.com/sp187/geekbang/final/internal/framework/middleware/trace"
)

var (
	ErrUserNotFound = &errs.Error{
		HttpCode: 200,
		Code:     600,
		Msg:      "未找到该用户",
	}
)

// UserService 定义用户服务， cmd中的使用用户服务时，需要初始化该对象，且同一个对象只初始化一次。
// 与用户所有相关的业务逻辑都依赖于该对象
type UserService struct {
	user UserRepo // 存储仓库
	// order OrderRepo
}

// UserRepo Repo 定义为了实现业务逻辑需要的与数据存储有关的接口，data目录下放实现该接口的不同存储库。比如postgresql和内存。
// 谁使用谁定义的原则，需要什么就定义什么。
type UserRepo interface {
	GetById(context.Context, uint) (User, error)
	Add(context.Context, User) error
}

// User 业务使用的领域对象DO(domain object)，用户业务内部逻辑围绕该对象展开。
type User struct {
	Id      int
	Name    string
	Gender  string // unknown male female
	Age     int
	Balance float64 // 用户余额
}

// NewUserService 新建一个用户服务，供cmd/下的web应用使用
func NewUserService(repo UserRepo) (*UserService, error) {
	us := &UserService{repo}
	return us, nil
}

// 业务相关的逻辑函数实现

// GetUser 查询用户
func (us *UserService) GetUser(ctx context.Context, id uint) (User, error) {
	// 链路追踪相关示例
	ctx, span := trace.GetTracer().Start(ctx, "GetUser")
	defer span.End()
	// 在存储中查找用户
	return us.user.GetById(ctx, id)
}

// AddUser 添加用户
func (us *UserService) AddUser(ctx context.Context, user User) error {
	// 往存储中添加用户
	return us.user.Add(ctx, user)
}
