package svc

import (
	"SAI-IM/apps/im/rpc/imclient"
	"SAI-IM/apps/social/api/internal/config"
	"SAI-IM/apps/social/api/internal/middleware"
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis

	socialclient.Social
	userclient.User
	imclient.Im

	IdempotenceMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		Redis:                 redis.MustNewRedis(c.Redisx),
		Social:                socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:                  userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:                    imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware(),
	}
}
