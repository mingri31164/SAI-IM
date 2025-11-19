package msgTransfer

import (
	"SAI-IM/apps/im/ws/ws"
	"SAI-IM/apps/task/mq/internal/svc"
	"SAI-IM/apps/task/mq/mq"
	"SAI-IM/pkg/bitmap"
	"SAI-IM/pkg/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"sync"
	"time"
)

// 默认值
var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota //默认不开启已读消息缓存合并处理
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer

	// 群聊消息缓存合并
	cache.Cache
	mu        sync.Mutex
	groupMsgs map[string]*groupMsgRead // 考虑到要处理多个群，用map存储
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}
	// 如果开启 已读消息缓存合并处理
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		// 最大计数
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			// 设置值
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}
		// 超时时间
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	//✨注意要协程运行
	go m.transfer()

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Infof("MsgReadTransfer ", value)
	var (
		data mq.MsgMarkRead
		//ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理--更新消息为已读
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	// 返回给接收者的结果map[string][已读记录]
	//✨注：因为已读记录需要通过 mq -> websocket -> 接收者，
	//      这个过程中避免类型问题，已读记录以string的方式传参会更合适

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	}
	// 判断消息类型
	switch data.ChatType {
	case constants.SingleChatType:
		// 直接推送
		m.push <- push
	case constants.GroupChatType:
		// 判断是否采用合并发送
		// 若不开启
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
			break
		}
		// 开启
		m.mu.Lock()
		defer m.mu.Unlock()
		push.SendId = "" // 因为群聊消息此时是合并的，不需要发送者ID
		// 已经有消息
		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			// 和并请求
			m.Infof("merge push: %v", push.ConversationId)
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			// 没有记录，创建
			m.Infof("create merge push %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}

	return nil
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	result := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读消息
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 设置当前发送者用户为已读状态
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		// 将已读消息（二进制）进行base64编码转换，这样可以保证在网络传输过程中的一个精度
		// 前端也可以使用base64进行解码，将已读消息还原为二进制
		result[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err := m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// 异步处理消息发送
func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("transfer err: %s", err.Error())
			}
		}

		if push.ChatType == constants.SingleChatType {
			// 类型有问题，不处理该消息
			continue
		}
		// 不采用合并推送
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据
		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].Clear()
			delete(m.groupMsgs, push.ConversationId)
		}
		m.mu.Unlock()
	}
}
