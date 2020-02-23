package rtmapi

import "encoding/json"

// WebSocketReply is passed from slack as a reply to client message
type WebSocketReply struct {
	OK *bool `json:"ok"`
	ReplyTo uint `json:"reply_to"`
	TimeStamp TimeStamp `json:"ts"`
	Text string `json:"text"`
}

// DecodeReply parses given reply payload from slack
func DecodeReply(input json.RawMessage) (*WebSocketReply, error) {
	reply := &WebSocketReply{}
	if err := json.Unmarshal(input, reply); err != nil {
		return nil, NewPayloadError(err.Error())
	}

	if reply.OK == nil {
		return nil, NewPayloadError("ok field is missing " + string(input))
	}

	return reply, nil
}