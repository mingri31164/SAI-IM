package main

import (
	"SAI-IM/apps/task/mq/internal/config"
	"SAI-IM/apps/task/mq/internal/handler"
	"SAI-IM/apps/task/mq/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	// 获取监听的ServiceContext服务
	listen := handler.NewListen(ctx)

	// 获取go-zero中的服务组对象
	serviceGroup := service.NewServiceGroup()
	//✨从listen中把服务信息加入到服务组中，由服务组统一管理所有的服务
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	// 通过服务组启动注册的所有服务
	serviceGroup.Start()
}
