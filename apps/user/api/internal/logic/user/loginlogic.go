package user

import (
	"SAI-IM/apps/user/rpc/user"
	"SAI-IM/pkg/constants"
	"context"
	"github.com/jinzhu/copier"

	"SAI-IM/apps/user/api/internal/svc"
	"SAI-IM/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 将用户ID和在线状态"1"存储到Redis的hash中，标记用户为在线。
	// 这里使用Redis来管理在线用户，是因为Redis的高并发读写性能和键值对存储特性适合此类场景。
	err = l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")
	if err != nil {
		// 如果设置Redis中用户在线状态失败，返回错误。
		return nil, err
	}
	return &res, nil
}
