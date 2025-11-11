package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex

	routes         map[string]HandlerFunc
	addr           string
	connToUser     map[*websocket.Conn]string
	userToConn     map[string]*websocket.Conn
	upgrader       websocket.Upgrader
	authentication Authentication
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),
		upgrader:       websocket.Upgrader{},
		authentication: new(authentication),
		Logger:         logx.WithContext(context.Background()),
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
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return
	}

	// 根据连接对象获取请求，根据请求查找路由并执行
	go s.handlerConn(conn)

}

func (s *Server) addConn(conn *websocket.Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *websocket.Conn) {

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 根据请求的method分发路由并执行
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法 %v，请检查",
				message.Method)))
		}
	}
}

/*
这里的 ...string 表示：
这个函数的参数 sendIds 是一个“可变数量的 string 参数”。
*/
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	/*
		Go 里的 ... 在调用函数时有一个反向功能：
		如果一个函数接收可变参数，而你已经有一个切片（比如 []string），
		想把切片里的每个元素当作独立参数传进去，就需要加上 ... 展开操作符。
	*/
	return s.Send(msg, s.GetConns(sendIds...)...)
}

// 根据连接对象执行任务处理
func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}

func (s *Server) Close(conn *websocket.Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	err := conn.Close()
	if err != nil {
		return
	}
}
