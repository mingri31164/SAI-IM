package friend

import (
	"SAI-IM/apps/social/rpc/social"
	"SAI-IM/pkg/constants"
	"SAI-IM/pkg/ctxdata"
	"context"

	"SAI-IM/apps/social/api/internal/svc"
	"SAI-IM/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// 从上下文中获取当前用户ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用服务层接口获取当前用户的所有好友列表
	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		// 如果获取好友列表失败，返回空的响应和错误信息
		return &types.FriendsOnlineResp{}, err
	}

	// 如果当前用户没有好友，直接返回空的在线状态响应
	if len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, nil
	}

	// 提取好友ID列表
	uids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		uids = append(uids, friend.UserId)
	}

	// 查询Redis缓存中的在线用户
	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		// 如果查询Redis缓存失败，返回错误信息
		return nil, err
	}

	// 构建在线状态的映射表
	resOnlineList := make(map[string]bool, len(uids))
	for _, uid := range uids {
		// 检查每个好友ID是否在在线用户列表中
		if _, ok := onlines[uid]; ok {
			// 如果好友在在线用户列表中，则标记为在线（即未找到，说明是离线状态）
			resOnlineList[uid] = true
		} else {
			// 如果好友不在在线用户列表中，则标记为离线
			resOnlineList[uid] = false
		}
	}

	// 返回好友在线状态的响应
	return &types.FriendsOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
