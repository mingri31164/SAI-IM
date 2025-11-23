package job

import (
	"context"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestWithRetry(t *testing.T) {
	var (
		ErrTest = errors.New("测试异常")
		// 默认每隔 1s 执行一次
		handler = func(ctx context.Context) error {
			t.Log("执行handler")
			return ErrTest
		}
	)
	type args struct {
		ctx     context.Context
		handler func(context.Context) error
		opts    []RetryOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		//✅ 测试用例 1：使用默认配置 → 应该超时退出
		{
			"1", args{
				ctx:     context.Background(),
				handler: handler,
				opts:    []RetryOptions{},
			}, ErrJobTimeout,
		},
		//✅ 测试用例 2：自定义重试超时 + 自定义延迟策略 → 最后返回 handler 原始错误
		{
			"2", args{
				ctx:     context.Background(),
				handler: handler,
				opts: []RetryOptions{
					WithRetryTimeout(3 * time.Second),
					WithRetryJetLagFunc(func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
						return 500 * time.Millisecond // 500ms未超时，抛出的不是超时错误，而是任务错误
					}),
				},
			}, ErrTest,
		},
		//✅ 测试用例 3：自定义 RetryFunc → 不进行重试，第一次失败后立即结束
		{
			"3", args{
				ctx:     context.Background(),
				handler: handler,
				opts: []RetryOptions{
					WithRetryFunc(func(ctx context.Context, retryCount int, err error) bool {
						return false
					}),
				},
			}, ErrTest,
		},
	}
	// 遍历执行定义的测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WithRetry(tt.args.ctx, tt.args.handler, tt.args.opts...); err != tt.wantErr {
				t.Errorf("WithRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
