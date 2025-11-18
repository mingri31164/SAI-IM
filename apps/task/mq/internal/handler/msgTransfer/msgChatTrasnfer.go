package msgTransfer

import (
	"SAI-IM/apps/im/immodels"
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/apps/task/mq/internal/svc"
	"SAI-IM/apps/task/mq/mq"
	"context"
	"encoding/json"
	"fmt"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

// Consume ✨实现 kq -> queue 中的Consume接口
// 不同于v1.1.8, 在v1.2.2中Consume接口的参数中新增了一个参数ctx
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, " value : ", value)
	var (
		data mq.MsgChatTransfer
		//ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	// 推送
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	// 更新会话
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
