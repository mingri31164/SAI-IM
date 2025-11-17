package msgTransfer

import (
	immodels "SAI-IM/apps/im/immodels"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/apps/task/mq/internal/svc"
	"SAI-IM/apps/task/mq/mq"
	"SAI-IM/pkg/constants"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
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

	switch data.ChatType {
	case constants.SingleChatType:
		return m.single(&data)
	case constants.GroupChatType:
		return m.group(ctx, &data)
	}

	return nil
}

func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, //此处FromId定义为系统角色，用于内部mq消息传递
		Data:      data,
	})
}

func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 获取群用户id
	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, 0, len(users.List))
	//fmt.Printf("group user: %+v", users.List)
	for _, members := range users.List {
		// 过滤掉发送者自己，避免重复给自己推送一次
		if members.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, members.UserId)
	}
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, //此处FromId定义为系统角色，用于内部mq消息传递
		Data:      data,
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
	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	// 更新会话
	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
