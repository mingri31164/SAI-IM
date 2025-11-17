package logic

import (
	"SAI-IM/apps/im/immodels"
	"SAI-IM/pkg/constants"
	"SAI-IM/pkg/xerr"
	"context"
	"github.com/pkg/errors"

	"SAI-IM/apps/im/rpc/im"
	"SAI-IM/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "conversationModel.FindByUserId err %v,req %v", err, in.UserId)
	}
	// 存在会话，但是会话列表可能为空，需要初始化
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}
	for s, conversation := range in.ConversationList {
		var oldTotal int
		if data.ConversationList[s] != nil {
			// 用户在当前会话中在此之前已读取的量
			oldTotal = data.ConversationList[s].Total
		}
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			// ✨用户传递的已读取的量+原本已读取的量 = 实际用户已读取量
			Total: int(conversation.Read) + oldTotal,
			Seq:   conversation.Seq,
		}
	}
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "conversationModel.Update err %v,req %v", err, data)
	}
	return &im.PutConversationsResp{}, nil
}
