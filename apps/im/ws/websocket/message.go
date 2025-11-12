package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0 // 用户消息
	FramePing FrameType = 0x1 // 心跳消息

	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// msg , id, seq
type Message struct {
	FrameType `json:"frameType"`

	Method string      `json:"method"`
	FormId string      `json:"formId"`
	Data   interface{} `json:"data"` // map[string]interface{}
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}
