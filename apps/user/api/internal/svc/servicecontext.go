package svc

import (
	"SAI-IM/apps/user/api/internal/config"
	"SAI-IM/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	// N * client =》 别名
)

// google.golang.org\grpc\server_config.go -> jsonMc jsonRetryPolicy
// grpc重试机制配置文件
//var retryPolicy = `{
//	"methodConfig" : [{
//		"name": [{
//			"service": "user.User"
//		}]
//		"waitForReady": true,
//		"retryPolicy": {
//			"maxAttempts": 5,
//			"initialBackoff": "0.001s",
//			"maxBackoff": "0.002s",
//			"backoffMultiplier": 1.0,
//			"retryableStatusCodes": ["UNKNOWN"]
//		}
//	}]
//}`

type ServiceContext struct {
	Config config.Config
	userclient.User
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Redis:  redis.MustNewRedis(c.Redisx),

		// 实践grpc重试机制
		//User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc,
		//	zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)))),
	}
}
