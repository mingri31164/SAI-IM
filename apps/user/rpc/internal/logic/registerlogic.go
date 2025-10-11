package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"SAI-IM/apps/user/rpc/internal/svc"
	"SAI-IM/apps/user/rpc/user"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// TODO 用户注册，添加默认密码，注册时无需输入密码

	return &user.RegisterResp{}, nil
}
