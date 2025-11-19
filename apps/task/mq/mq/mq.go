package mq

import "SAI-IM/pkg/constants"

// 定义Kafka接收消息的格式
type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string   `json:"sendId"`
	RecvId             string   `json:"recvId"`
	RecvIds            []string `mapstructure:"recvIds"` // 群聊消息推送时考虑多个用户的接收
	SendTime           int64    `json:"sendTime"`

	constants.MType `json:"mType"`
	Content         string `json:"content"`
}

// 消息已读未读
type MsgMarkRead struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string   `json:"sendId"`
	RecvId             string   `json:"recvId"`
	MsgIds             []string `json:"msgIds"`
}
