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

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新会话
func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (resp *types.PutConversationsResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	var conversationList = make(map[string]*im.Conversation)
	copier.Copy(&conversationList, req.ConversationList)

	_, err = l.svcCtx.PutConversations(l.ctx, &im.PutConversationsReq{
		UserId:           uid,
		ConversationList: conversationList,
	})
	if err != nil {
		return nil, err
	}

	return &types.PutConversationsResp{}, nil
}
