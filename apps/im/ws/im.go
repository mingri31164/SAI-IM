package main

import (
	"SAI-IM/apps/im/ws/internal/config"
	"SAI-IM/apps/im/ws/internal/svc"
	"SAI-IM/apps/im/ws/websocket"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svc.NewServiceContext(c)

	srv := websocket.NewServer(c.ListenOn)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
