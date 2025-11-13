package conversation

import (
	"SAI-IM/apps/im/ws/internal/logic"
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/pkg/constants"
	"context"
	"github.com/mitchellh/mapstructure"
	"time"
)

// 单独定义会话是因为聊天本身是需要在会话中去完成的，所以不在用户中去创建会话
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		// 对map[string]interface{}类型进行转换
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		switch data.ChatType {
		case constants.SingleChatType:
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}
			srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				SendTime:       time.Now().UnixMilli(),
				Msg:            data.Msg,
			}), data.RecvId)
		}
	}
}
