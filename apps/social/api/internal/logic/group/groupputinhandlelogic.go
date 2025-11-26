package group

import (
	"SAI-IM/apps/im/rpc/imclient"
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/pkg/constants"
	"SAI-IM/pkg/ctxdata"
	"context"
	"strconv"

	"SAI-IM/apps/social/api/internal/svc"
	"SAI-IM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleReq) (resp *types.GroupPutInHandleResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	// 转换GroupReqId从string到int32
	groupReqId, err := strconv.ParseInt(req.GroupReqId, 10, 32)
	if err != nil {
		return nil, err
	}

	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   int32(groupReqId),
		GroupId:      req.GroupId,
		HandleUid:    uid,
		HandleResult: req.HandleResult,
	})

	if res.GroupId == "" {
		return
	}
	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	_, err = l.svcCtx.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return nil, err
}
