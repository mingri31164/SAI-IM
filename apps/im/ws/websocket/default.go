package websocket

import (
	"math"
	"time"
)

const (
	// 默认最大空闲时间，当前项目中默认没有这项检测，设置为最大int类型
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)
	// ACK默认超时时间
	defaultAckTimeout = 30 * time.Second
	// 定义websocket并发的默认量级
	defaultConCurrency = 30
)
