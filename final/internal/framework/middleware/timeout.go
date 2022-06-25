package mid

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type timeKey int

// TimeControl 用于请求的超时控制
func TimeControl(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	now := time.Now()
	ctx := r.Context()
	ctx = context.WithValue(ctx, timeKey(1), now)
	ttl := r.Header.Get("X-Time-To-Live")
	if ttl != "" {
		ms, err := strconv.ParseInt(ttl, 10, 64)
		if err == nil {
			if ms < 3 {
				// 减去3毫秒内网传输时间
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			}
			ctx, _ = context.WithTimeout(ctx, time.Duration(ms-3)*time.Millisecond)
		}
	}
	r = r.WithContext(ctx)
	next(w, r)
}

func TimeFromReq(r *http.Request) time.Time {
	t, ok := r.Context().Value(timeKey(1)).(time.Time)
	if ok {
		return t
	} else {
		return time.Now()
	}
}

func TimeFromCtx(ctx context.Context) time.Time {
	t, ok := ctx.Value(timeKey(1)).(time.Time)
	if ok {
		return t
	} else {
		return time.Now()
	}
}
