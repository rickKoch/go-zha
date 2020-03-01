package rtmapi

import "sync"

// OutgoingEvent is empty interface that marks outgoing events.
type OutgoingEvent interface{}

// OutgoingCommonEvent takes care of some common fields
type OutgoingCommonEvent struct {
	CommonEvent
	ID uint `json:"id"`
}

// OutgoingMessage represents a simple message sent from client to Slack
type OutgoingMessage struct {
	OutgoingCommonEvent
	ID      uint   `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// NewOutgoingMessage creates new OutgoingMessage instace
func NewOutgoingMessage(eventID *OutgoingEventID, message *TextMessage) *OutgoingMessage {
	return &OutgoingMessage{
		Channel: message.channel,
		Text:    message.text,
		OutgoingCommonEvent: OutgoingCommonEvent{
			ID:          eventID.Next(),
			CommonEvent: CommonEvent{Type: MESSAGE},
		},
	}
}

// Ping struct
type Ping struct {
	OutgoingCommonEvent
}

// NewPing creates new Ping instance
func NewPing(eventID *OutgoingEventID) *Ping {
	return &Ping{
		OutgoingCommonEvent: OutgoingCommonEvent{
			ID:          eventID.Next(),
			CommonEvent: CommonEvent{Type: PING},
		},
	}
}

// OutgoingEventID struct
type OutgoingEventID struct {
	id    uint
	mutex *sync.Mutex
}

// Next generetes the next id
func (m *OutgoingEventID) Next() uint {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.id++
	return m.id
}

// NewOutgoingEventID creates new  OutgoingEventID
func NewOutgoingEventID() *OutgoingEventID {
	return &OutgoingEventID{
		id:    0,
		mutex: &sync.Mutex{},
	}
}
