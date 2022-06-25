package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server interface {
	Route(string, ServiceFunc, ...string) // 添加路由
	Start(address string) error           // 启动服务，监听对应的端口
	Append(handler Handler)               // 添加调用链  request -> Handler 1 -> Handler 2 -> ... -> Handler N -> response
	Shutdown(ctx context.Context) error   // 服务退出
}

type ServiceFunc func(c *Context)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(w, r, next)
}

// WarpHandler 将原生http.Handler转为Handler
func WarpHandler(handler http.Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	})
}

// WarpServiceFunc 将http.Handler转为ServiceFunc
func WarpServiceFunc(handler http.Handler) ServiceFunc {
	return func(c *Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Server的一个实现
type httpServer struct {
	Name    string
	router  *mux.Router
	handler *negroni.Negroni
	server  *http.Server
}

func (h *httpServer) Route(s string, f ServiceFunc, methods ...string) {
	var r *mux.Route
	if strings.HasSuffix(s, "/*") {
		s = s[:len(s)-1]
		r = h.router.PathPrefix(s).Handler(http.StripPrefix(s, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			c := NewContext(w, req)
			c.vars = mux.Vars(req)
			f(c)
		})))
	} else {
		r = h.router.HandleFunc(s, func(w http.ResponseWriter, req *http.Request) {
			c := NewContext(w, req)
			c.vars = mux.Vars(req)
			f(c)
		})
	}
	r.Methods(methods...)

}

func (h *httpServer) Start(address string) error {
	h.Append(WarpHandler(h.router))
	h.server.Addr = address
	h.server.Handler = h.handler
	return h.server.ListenAndServe()
}

func (h *httpServer) Append(handler Handler) {
	h.handler.Use(handler)
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

// NewHttpServer 新建一个http server
func NewHttpServer(name string) *httpServer {
	return &httpServer{name, mux.NewRouter(), negroni.New(), &http.Server{}}
}
