package handler

import (
	"SAI-IM/apps/im/ws/internal/handler/conversation"
	"SAI-IM/apps/im/ws/internal/handler/user"
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
)

// 与api层不同，这里注册的是我们自定义的websocket服务
func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
	})
}
