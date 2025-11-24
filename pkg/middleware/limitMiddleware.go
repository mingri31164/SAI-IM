package middleware

import (
	"SAI-IM/pkg/constants"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type LimitMiddleware struct {
	redisCfg redis.RedisConf
	*limit.TokenLimiter
}

func NewLimitMiddleware(redisCfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{
		redisCfg: redisCfg,
	}
}

// TokenLimitHandler @Param 限流速率(每秒生成多少令牌), 桶容量
func (m *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	m.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(m.redisCfg), constants.RedisTokenLimitKey)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 如果获取令牌成功，继续执行业务
			if m.TokenLimiter.AllowCtx(r.Context()) {
				next(w, r)
				return
			}
			// 限流（抛出指定状态码，go-zero限流源码自动识别并执行限流）
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
}
