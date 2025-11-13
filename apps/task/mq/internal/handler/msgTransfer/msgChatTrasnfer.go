package msgTransfer

import (
	"SAI-IM/apps/task/mq/internal/svc"
	"context"
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
// 不同于v1.1.8, v1.2.2中，Consume接口的参数中新增了一个参数ctx
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, " value : ", value)
	return nil
}
