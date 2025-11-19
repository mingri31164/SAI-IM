package constants

type MType int

const (
	TextMType MType = iota
)

// ChatType 聊天类型
type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

// ContentType 内容类型
type ContentType int

const (
	ContentChatMsg  ContentType = iota // 聊天
	ContentMakeRead                    // 已读记录
)
