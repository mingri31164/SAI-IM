package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}

type authentication struct{}

func (*authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	// 暂定为默认通过，因为websocket连接是在客户端登录鉴权后才开始的
	return true
}
func (*authentication) UserId(r *http.Request) string {
	// 请求id: 1.userId  2.时间戳
	// userId在url拼接中传递
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
