package svc

import (
	immodels "SAI-IM/apps/im/immodels"
	"SAI-IM/apps/im/ws/websocket"
	"SAI-IM/apps/social/rpc/socialclient"
	"SAI-IM/apps/task/mq/internal/config"
	"SAI-IM/pkg/constants"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
)

type ServiceContext struct {
	config.Config
	WsClient websocket.Client
	*redis.Redis

	// 对数据库MongoDB的连接
	immodels.ChatLogModel
	// 会话模型
	immodels.ConversationModel
	// 提供社交服务的rpc支持
	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
