package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	xerrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var stopSignal = make(chan bool)

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	server := NewHttpServer("test")
	server.Route("/", HelloWorld)
	server.Route("/stop", StopServer)

	// 1
	eg.Go(func() error {
		fmt.Println("listen port 8080")
		err := server.Serve(":8080")
		fmt.Printf("exit goroutine 1: %+v\n", err)
		return err
	})

	// 2
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("exit goroutine 2: some errors happened in goroutine 1 or 3")
		case <-stopSignal:
			fmt.Println("exit goroutine 2: receive server stop signal")
		}
		return server.Stop()
	})

	// 3
	eg.Go(func() error {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
			fmt.Println("exit goroutine 3: some errors happened in goroutine 1 or 2")
		case s := <-sig:
			err = xerrors.Errorf("exit goroutine 3: system signal call %+v", s)
			fmt.Println(err.Error())
		}
		return err
	})

	fmt.Printf("server exit: %+v\n", eg.Wait())
}

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}

func StopServer(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("goodbye"))
	close(stopSignal)
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
