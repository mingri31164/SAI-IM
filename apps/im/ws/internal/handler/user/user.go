package user

import (
	"SAI-IM/apps/im/ws/internal/svc"
	websocketx "SAI-IM/apps/im/ws/websocket"
	"github.com/gorilla/websocket"
)

func OnLine(srv *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		ids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(u[0], ids), conn)
		srv.Info("err", err)
	}
}
