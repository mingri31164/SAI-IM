package mqclient

import (
	"SAI-IM/apps/task/mq/mq"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

// 定义对于聊天消息转化的会话客户端，提供给websocket服务进行使用
type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	//✨使用go-zero提供的第三方库中，kafka的pusher对象来完成消息的发送
	pusher *kq.Pusher
}

// opts ...kq.PushOption是对kafka扩展组件的设置
func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(context.Background(), string(body))
}
