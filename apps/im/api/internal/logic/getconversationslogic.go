package logic

import (
	"SAI-IM/apps/im/rpc/im"
	"SAI-IM/pkg/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"SAI-IM/apps/im/api/internal/svc"
	"SAI-IM/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话
func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	data, err := l.svcCtx.GetConversations(l.ctx, &im.GetConversationsReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	var conversationList = make(map[string]*types.Conversation)
	copier.Copy(&conversationList, data.ConversationList)

	return &types.GetConversationsResp{ConversationList: conversationList}, nil
}
