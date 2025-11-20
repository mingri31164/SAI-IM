package main

import (
	"SAI-IM/pkg/configserver"
	"SAI-IM/pkg/resultx"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"

	"SAI-IM/apps/user/api/internal/config"
	"SAI-IM/apps/user/api/internal/handler"
	"SAI-IM/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// go-zero配置默认加载方式
	//conf.MustLoad(*configFile, &c)

	var configs = "user-api.yaml"

	// 配置中心加载方式
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "118.178.120.11:3379",
		ProjectKey:     "3c46a0407be60a1f00731ab8e9575df2",
		Namespace:      "user",
		Configs:        configs,
		ConfigFilePath: "../etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c)
	if err != nil {
		panic(err)
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
