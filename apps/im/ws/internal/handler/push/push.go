package push

import (
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/pkg/constants"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			_ = srv.Send(websocket.NewErrMessage(err))
			return
		}

		// 发送目标
		// todo: 这里的recv是一个切片
		switch data.ChatType {
		case constants.SingleChatType:
			err := single(srv, &data, data.RecvId)
			if err != nil {
				srv.Error(err)
			}
		case constants.GroupChatType:
			err := group(srv, &data)
			if err != nil {
				srv.Error(err)
			}
		default:

		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
		return nil
	}
	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
			ReadRecords: data.ReadRecords,
		},
	}), rconn)

}

func group(srv *websocket.Server, data *ws.Push) (err error) {
	//fmt.Println("group push")
	// ✨此处群聊实现是在私聊的基础上进行迭代
	for _, id := range data.RecvIds {
		func(id string) {
			// 此处Schedule为threading.TaskRunner下的并发调用方法
			srv.Schedule(func() {
				err = single(srv, data, id)
			})
		}(id)
		//fmt.Println(id)
	}
	return
}
