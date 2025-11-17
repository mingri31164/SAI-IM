package conversation

import (
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/apps/task/mq/mq"
	"SAI-IM/pkg/constants"
	"SAI-IM/pkg/wuid"
	"github.com/mitchellh/mapstructure"
	"time"
)

// Chat 单独定义会话是因为聊天本身是需要在会话中去完成的，所以不在用户中去创建会话
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		// 对map[string]interface{}类型进行转换
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		// 如果没有传递会话id，则根据聊天类型分别创建会话id
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				// 群的会话id就是接收者（群id）
				data.ConversationId = data.RecvId
			default:
			}
		}
		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixMilli(),
			MType:          data.MType,
			Content:        data.Content,
		})
		if err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
