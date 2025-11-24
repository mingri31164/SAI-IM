package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// SyncLimiterInterceptor 底层基于channel管道阻塞实现并发限流
func SyncLimiterInterceptor(maxCount int) grpc.UnaryServerInterceptor {
	l := syncx.NewLimit(maxCount) // 定义最大并发数量
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if l.TryBorrow() { // 判断当前是否进行了限流
			defer func() {
				if err := l.Return(); err != nil {
					logx.Error(err)
				}
			}()
		} else {
			logx.Errorf("concurrent connections over %d,rejected with code %d",
				maxCount, http.StatusServiceUnavailable)
			return nil, status.Error(codes.Unavailable, "concurrent connection over limit")
		}
		resp, err = handler(ctx, req)
		return
	}
}
