package interceptor

import (
	"SAI-IM/pkg/xerr"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Idempotent interface {
	// Identify 获取请求的标识
	Identify(ctx context.Context, method string) string
	// IsIdempotentMethod 是够支持幂等性
	IsIdempotentMethod(fullMethod string) bool
	// TryAcquire 幂等性的验证
	TryAcquire(ctx context.Context, id string) (resp any, isAcquire bool)
	// SaveResp 执行之后结果的保存
	SaveResp(ctx context.Context, id string, resp any, respErr error) error
}

var (
	// TKey 请求任务标识
	TKey = "sai-im-idempotence-task-id"
	// DKey 设置rpc调度中的rpc请求的标识
	DKey = "sai-im-idempotence-dispatch-key"
)

// ContextWithVal 添加到上下文方便客户端获取
func ContextWithVal(ctx context.Context) context.Context {
	// 设置请求id
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

// NewIdempotenceClient grpc客户端的拦截器
func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取唯一的key
		identify := idempotent.Identify(ctx, method)
		// 在rpc请求中设置头部信息
		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: []string{identify},
		})
		// 请求
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewIdempotenceServer grpc服务端的拦截器
func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 获取请求id
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		if len(identify) == 0 || !idempotent.IsIdempotentMethod(info.FullMethod) {
			// 不进行幂等处理
			return handler(ctx, req)
		}

		fmt.Println("----", "请求进行幂等处理", identify)
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			resp, err = handler(ctx, req)
			fmt.Println("---- 执行任务")
			// 保存执行之后的结果
			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}
			return resp, err
		}
		// 任务已经执行完了
		if r != nil {
			fmt.Println("----", "任务已经执行完了")
			return r, nil
		}
		// 任务还在执行
		//🔥注意：因为需要grpc的重试，所以此处需要使用grpc的错误码
		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行"+
			"id %v", identify[0])))
	}
}
