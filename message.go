package zha

import (
	"context"
	"fmt"
)

// Message struct
type Message struct {
	Context  context.Context
	Text     string
	ChannelD string
	Matches  []string

	adapter Adapter
}

// Respond test
func (msg *Message) Respond(text string, args ...interface{}) {
	if len(args) > 0 {
		text = fmt.Sprintf(text, args...)
	}

	_ = msg.adapter.Send(text, msg.ChannelD)
}
