package service

import (
	"fmt"
	"github.com/sp187/geekbang/final/cmd/temp-app/service/dto"
	"strconv"
	"time"

	"github.com/sp187/geekbang/final/internal/framework"
	"github.com/sp187/geekbang/final/internal/framework/error"
	"github.com/sp187/geekbang/final/internal/framework/middleware/trace"
	"github.com/sp187/geekbang/final/internal/framework/web"
	utils "gitlab.bj.sensetime.com/sense-remote/product/sense-layers/rs-component-utils.git"
)

func SayHello(c *web.Context) {
	var name string
	name = c.Var("name")
	fw.GetLogger().Info("欢迎光临 %s", name)
	c.OKJsonResponse(struct {
		Message string `json:"message"`
	}{fmt.Sprintf("欢迎光临 %s", name)})
}

func GetUserByID(c *web.Context) {
	// 添加链路追踪信息示例
	ctx, span := trace.GetTracer().Start(c.Context(), "GetUserByID")
	defer span.End()

	id := c.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.BadJsonResponse(errs.BadRequest)
	}
	user, err := userService.GetUser(ctx, uint(idInt))
	if err != nil {
		c.BadJsonResponse(err)
		return
	}
	resp := dto.ToUserResp(&user)
	c.OKJsonResponse(resp)
	return
}

func AddUser(c *web.Context) {
	user := dto.UserDTO{}
	err := c.ReadJson(&user)
	if err != nil {
		c.BadJsonResponse(utils.ErrorInvalidParameter)
		return
	}
	// userService是一个服务对象
	err = userService.AddUser(c.Context(), user.ToDO())
	if err != nil {
		c.BadJsonResponse(err)
		return
	}
	c.OKJsonResponse(nil)
	return
}

func SlowAPI(c *web.Context) {
	time.Sleep(10 * time.Second)
	c.OKJsonResponse(nil)
}
