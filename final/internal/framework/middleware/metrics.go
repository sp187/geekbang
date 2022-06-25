package mid

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
)

type SimpleMetrics struct {
	inFlightReq int64 // 处理中的请求数
	totalReq    int64
	// todo 其他统计指标、参数
}

// NewMetrics 创建指标结构体，不要与GracefulShutdown同时使用
func NewMetrics() *SimpleMetrics {
	return &SimpleMetrics{}
}

// CountRequestHandler 记录正在执行的请求数
func (s *SimpleMetrics) CountRequestHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	atomic.AddInt64(&s.inFlightReq, 1)
	atomic.AddInt64(&s.totalReq, 1)
	next(w, r)
	atomic.AddInt64(&s.inFlightReq, -1)
}

// GetInFlightRequest 获取处理中的请求数
func (s *SimpleMetrics) GetInFlightRequest() int64 {
	return atomic.LoadInt64(&s.inFlightReq)
}

func (s *SimpleMetrics) GetTotalRequest() int64 {
	return atomic.LoadInt64(&s.totalReq)
}

type MetricsWithShutdown struct {
	SimpleMetrics
	close int32         // 服务关闭标志，大于0表示关闭
	done  chan struct{} // 所有请求处理完毕后往该通道发送数据
}

// NewMetricsWithShutdown 创建一个可优雅退出的结构体，不要与SimpleMetrics同时使用
func NewMetricsWithShutdown() *MetricsWithShutdown {
	return &MetricsWithShutdown{
		done: make(chan struct{}),
	}
}

// ShutdownHandler 用于优雅退出的handler
func (g *MetricsWithShutdown) ShutdownHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cl := atomic.LoadInt32(&g.close)
	if cl > 0 {
		// 拒绝所有的请求
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	atomic.AddInt64(&g.inFlightReq, 1)
	atomic.AddInt64(&g.totalReq, 1)

	next(w, r)
	cl = atomic.LoadInt32(&g.close)
	n := atomic.AddInt64(&g.inFlightReq, -1)
	// 已经开始关闭了，而且请求数为0
	if cl > 0 && n == 0 {
		g.done <- struct{}{}
	}
}

// RejectAndWaiting 将会拒绝新的请求，并且等待处理中的请求
func (g *MetricsWithShutdown) RejectAndWaiting(ctx context.Context) error {
	// 设置关闭标志
	atomic.AddInt32(&g.close, 1)
	if atomic.LoadInt64(&g.inFlightReq) == 0 {
		return nil
	}
	done := ctx.Done()
	select {
	case <-done:
		return errors.New("shutdown timeout")
	case <-g.done:
		// shutdown successfully
	}
	return nil
}
