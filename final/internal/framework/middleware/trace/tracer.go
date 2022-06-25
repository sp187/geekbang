package trace

import (
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	initTracerOnce sync.Once
	tp             *sdktrace.TracerProvider
)

// InitStdTracer 初始化一个标准控制台的TracerProvider，链路信息直接打印
func InitStdTracer() *sdktrace.TracerProvider {
	if tp != nil {
		return tp
	}
	initTracerOnce.Do(func() {
		var (
			err      error
			exporter *stdout.Exporter
		)
		exporter, err = stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			log.Fatal("InitStdTracer fail")
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		return
	})
	return tp
}

func GetTracer() trace.Tracer {
	if tp == nil {
		InitStdTracer()
	}
	return tp.Tracer("")
}
