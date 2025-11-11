package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type Server struct {
	addr     string
	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		upgrader: websocket.Upgrader{},
		Logger:   logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// 处理运行过程中可能会抛出的系统性panic，为避免服务崩溃，需要恢复并捕获异常
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	// 获取连接对象
	_, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return
	}

	// 根据连接对象获取请求

}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}
