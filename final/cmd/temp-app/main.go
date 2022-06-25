package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sp187/geekbang/final/cmd/temp-app/service"
	"github.com/sp187/geekbang/final/internal/framework"
	"github.com/sp187/geekbang/final/internal/framework/middleware"
	"github.com/sp187/geekbang/final/internal/framework/middleware/trace"
	"github.com/sp187/geekbang/final/internal/framework/web"
)

func main() {
	// log filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// loads environment values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// 初始化一些全局变量
	err := service.InitAll()
	if err != nil {
		panic("initialize environment fail:" + err.Error())
	}

	// 创建一个新的server
	server := web.NewHttpServer("test")

	// 可以创建一个admin server，用于高级权限的一些内部接口
	adminServer := web.NewHttpServer("admin")

	shutdown := mid.NewMetricsWithShutdown()

	// 添加日志中间件
	server.Append(mid.DefaultLogger())

	// 添加链路追踪的例子
	server.Append(web.HandlerFunc(trace.Handler("template")))
	tp := trace.InitStdTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	// 添加可记录请求数量的优雅退出中间件
	server.Append(web.HandlerFunc(shutdown.ShutdownHandler))

	// 路由
	// info路由可用来做心跳检测
	server.Route("/hello", service.SayHello, "POST", "GET")
	server.Route("/slow", service.SlowAPI, "GET")
	server.Route("/user", service.GetUserByID, "GET")
	server.Route("/user", service.AddUser, "POST")

	// 获取当前请求数和历史总请求数
	adminServer.Route("/stat", func(c *web.Context) {
		c.WriteJson(200, struct {
			Current int64 `json:"req"`
			Total   int64 `json:"total"`
		}{shutdown.GetInFlightRequest(), shutdown.GetTotalRequest()})
	})
	// 添加一个控制进程退出的api示例
	adminServer.Route("/shutdown", func(c *web.Context) {
		fw.GetLogger().Info("stopping server...")
		adminServer.Shutdown(context.TODO())
	})

	// 按给定端口启动服务
	port := fw.GetServicePort()
	adminPort := fw.GetAdminPort()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动普通服务
	go func() {
		fw.GetLogger().Info(fmt.Sprintf("start service on port %s", port))
		err := server.Start(":" + port)
		if err != nil {
			cancel()
		}
	}()

	// 启动管理服务
	go func() {
		fw.GetLogger().Info(fmt.Sprintf("start admin service on port %s", adminPort))
		err := adminServer.Start(":" + adminPort)
		if err != nil {
			cancel()
		}
	}()

	// 可接收系统信号优雅退出 任意server返回错误均会导致整个进程退出
	web.GracefulShutdown(ctx, shutdown.RejectAndWaiting, web.ShutdownServerHook(server, adminServer))

}
