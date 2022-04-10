package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/config"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var configPath string

var stopSignal = make(chan bool)

func init() {
	flag.StringVar(&configPath, "conf", "config/config.yaml", "-conf: config file path")
}

func main() {
	InitEnv()
	eg, ctx := errgroup.WithContext(context.Background())
	server := NewHttpServer("test")
	server.Route("/account", service.GetAccount)
	server.Route("/account/update", service.UpdateAccount)

	eg.Go(func() error {
		fmt.Println("listen port ", config.GetServicePort())
		err := server.Serve(":" + config.GetServicePort())
		fmt.Printf("exit goroutine 1: %+v\n", err)
		return err
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("exit goroutine 2: some errors happened in goroutine 1 or 3")
		case <-stopSignal:
			fmt.Println("exit goroutine 2: receive server stop signal")
		}
		return server.Stop()
	})

	eg.Go(func() error {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
			fmt.Println("exit goroutine 3: some errors happened in goroutine 1 or 2")
		case s := <-sig:
			err = errors.Errorf("exit goroutine 3: system signal call %+v", s)
			fmt.Println(err.Error())
		}
		return err
	})

	fmt.Printf("server exit: %+v\n", eg.Wait())
}

type Server interface {
	Route(pattern string, handlerFunc http.HandlerFunc)
	Serve(address string) error
	Stop() error
}

type httpServer struct {
	Name      string
	server    http.Server
	serverMux *http.ServeMux
}

func (hs *httpServer) Route(pattern string, handlerFunc http.HandlerFunc) {
	hs.serverMux.HandleFunc(pattern, handlerFunc)
}

func (hs *httpServer) Serve(address string) error {
	hs.server = http.Server{
		Addr:    address,
		Handler: hs.serverMux,
	}
	return hs.server.ListenAndServe()
}

func (hs *httpServer) Stop() error {
	return hs.server.Shutdown(context.TODO())
}

func NewHttpServer(name string) Server {
	return &httpServer{Name: name, serverMux: http.NewServeMux()}
}
