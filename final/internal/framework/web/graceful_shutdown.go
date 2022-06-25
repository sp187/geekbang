package web

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ShutdownSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL,
		syscall.SIGABRT, syscall.SIGTERM,
	}
)

// GracefulShutdown 优雅退出， hooks函数用于添加退出时需要执行的逻辑
func GracefulShutdown(ctx context.Context, hooks ...func(ctx context.Context) error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	select {
	case <-ctx.Done():
		fmt.Printf("shutdown with context cancel: %v", ctx.Err())
		os.Exit(0)
	case sig := <-signals:
		fmt.Printf("get signal %s, application will shutdown \n", sig)
		// 超时强制退出当前进程
		time.AfterFunc(5*time.Minute, func() {
			fmt.Printf("shutdown timeout, quit forced")
			os.Exit(1)
		})
		for _, f := range hooks {
			ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
			err := f(ctx)
			if err != nil {
				fmt.Printf("failed to run hook, err: %v \n", err)
			}
			cancel()
		}
		os.Exit(0)
	}
}

func ShutdownServerHook(servers ...Server) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		for _, server := range servers {
			server := server
			wg.Add(1)
			go func() {
				server.Shutdown(ctx)
				wg.Done()
			}()
		}
		done := make(chan struct{})
		go func() {
			wg.Wait()
			done <- struct{}{}
		}()

		select {
		case <-done:
			fmt.Println("all server shutdown")
			return nil
		case <-ctx.Done():
			return errors.New("server shutdown timeout")
		}
	}
}
