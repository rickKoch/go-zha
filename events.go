package zha

// InitEvent struct
type InitEvent struct{}

// ShutdownEvent struct
type ShutdownEvent struct{}

// ReciveMessageEvent struct
type ReciveMessageEvent struct {
	Text     string
	ChannelD string
}
