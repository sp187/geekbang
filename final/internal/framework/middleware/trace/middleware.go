package trace

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Handler 链路追踪handler，遵循OpenTelemetry的api，接收来自客户端的链路信息
func Handler(name string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		otelHandler := otelhttp.NewHandler(next, name)
		otelHandler.ServeHTTP(w, req)
	}
}
