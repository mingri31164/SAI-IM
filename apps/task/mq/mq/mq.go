package mq

import "SAI-IM/pkg/constants"

// 定义Kafka接收消息的格式
type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string `json:"sendId"`
	RecvId             string `json:"recvId"`
	SendTime           int64  `json:"sendTime"`

	constants.MType `json:"mType"`
	Content         string `json:"content"`
}
