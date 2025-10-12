package logic

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"SAI-IM/apps/user/rpc/internal/svc"
	"SAI-IM/apps/user/rpc/models"
	"SAI-IM/apps/user/rpc/user"
	"SAI-IM/pkg/ctxdata"
	"SAI-IM/pkg/encrypy"
	"SAI-IM/pkg/xerr"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	u := &models.User{}
	var err error
	// 1.检查用户是否存在(phone)
	err = l.svcCtx.CSvc.GetUserByPhone(u, in.Phone)
	if err != nil {
		if u.ID == "" {
			return nil, errors.WithStack(xerr.PhoneNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find api by phone "+
			"%v err %v ", in.Phone, err)
	}

	// 2. 密码验证
	if !encrypy.ValidatePasswordHash([]byte(u.Password), []byte(in.Password)) {
		return nil, errors.WithStack(xerr.UserPwdErr)
	}
	// 3. 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, u.ID)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "etxdata get jwt token"+
			" err %v ", in.Phone)
	}
	return &user.LoginResp{
		Id:     u.ID,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
