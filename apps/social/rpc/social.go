package main

import (
	"SAI-IM/pkg/interceptor"
	"SAI-IM/pkg/interceptor/rpcserver"
	"flag"
	"fmt"

	"SAI-IM/apps/social/rpc/internal/config"
	"SAI-IM/apps/social/rpc/internal/server"
	"SAI-IM/apps/social/rpc/internal/svc"
	"SAI-IM/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	//✨ 在客户端传递参数即可，但在服务端需要传递出Redis的配置信息
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
