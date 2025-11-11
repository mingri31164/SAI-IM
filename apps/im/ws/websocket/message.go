package websocket

// msg , id, seq
type Message struct {
	Method string      `json:"method"`
	FormId string      `json:"formId"`
	Data   interface{} `json:"data"` // map[string]interface{}
}
