package svc

import (
	"SAI-IM/apps/im/api/internal/config"
	"SAI-IM/apps/im/rpc/im"
	"SAI-IM/apps/im/rpc/imclient"
	"SAI-IM/apps/social/rpc/social"
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/apps/user/rpc/user"
	"SAI-IM/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	im.ImClient
	social.SocialClient
	user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ImClient:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		SocialClient: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		UserClient:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
