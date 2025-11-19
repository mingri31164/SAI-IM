package msgTransfer

import (
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/apps/task/mq/internal/svc"
	"SAI-IM/apps/task/mq/mq"
	"SAI-IM/pkg/bitmap"
	"SAI-IM/pkg/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

type MsgReadTransfer struct {
	*baseMsgTransfer
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Infof("MsgReadTransfer ", value)
	var (
		data mq.MsgMarkRead
		//ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理--更新消息为已读
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	// 返回给接收者的结果map[string][已读记录]
	//✨注：因为已读记录需要通过 mq -> websocket -> 接收者，
	//      这个过程中避免类型问题，已读记录以string的方式传参会更合适

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	})

}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	result := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读消息
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 设置当前发送者用户为已读状态
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		// 将已读消息（二进制）进行base64编码转换，这样可以保证在网络传输过程中的一个精度
		// 前端也可以使用base64进行解码，将已读消息还原为二进制
		result[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err := m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
