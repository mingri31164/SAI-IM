package msgTransfer

import (
	"SAI-IM/apps/im/ws/ws"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu sync.Mutex

	conversationId string // 会话ID

	push     *ws.Push      // 用于记录消息
	pushChan chan *ws.Push // 推送用户合并消息的通道

	count    int       // 计数
	pushTime time.Time // 上次推送时间

	done chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushChan chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushChan:       pushChan,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	return m
}
