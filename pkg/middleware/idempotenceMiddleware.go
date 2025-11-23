package middleware

import (
	"SAI-IM/pkg/interceptor"
	"net/http"
)

type IdempotenceMiddleware struct {
}

func NewIdempotenceMiddleware() *IdempotenceMiddleware {
	return &IdempotenceMiddleware{}
}

func (m *IdempotenceMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 调用拦截器生成请求唯一标识并存入上下文中
		r = r.WithContext(interceptor.ContextWithVal(r.Context()))
		next(w, r)
	}
}
