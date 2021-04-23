package middleware

import (
	"context"
	"goweb教程/controllers"
	"net/http"
	"time"
)

// BasicAuthMiddleware ..
type TimeoutMiddleware struct {
	Next http.Handler
}

func (tm *TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}

	ctx := r.Context()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second) // 修改 reqeust 的上下文 设置超时，返回新的ctx
	r.WithContext(ctx)                               //reqeust 使用新的ctx

	ch := make(chan struct{})
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done(): // 如果 timeout 时间到了
		controllers.ResponseJsonError(w, r, http.StatusRequestTimeout, "超时")
		return
	}
	tm.Next.ServeHTTP(w, r)

}
