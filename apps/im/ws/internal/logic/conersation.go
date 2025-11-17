package logic

import (
	"SAI-IM/apps/im/immodels"
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/pkg/wuid"
	"context"
	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

// 处理单聊业务，具体处理方式最终要返回给上层conversation中去使用
func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	time.Sleep(time.Minute)
	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
