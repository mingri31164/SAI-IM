package friend

import (
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/pkg/ctxdata"
	"context"
	"strconv"

	"SAI-IM/apps/social/api/internal/svc"
	"SAI-IM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// 转换FriendReqId从string到int32
	friendReqId, err := strconv.ParseInt(req.FriendReqId, 10, 32)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  int32(friendReqId),
		UserId:       ctxdata.GetUId(l.ctx),
		HandleResult: req.HandleResult,
	})

	return
}
