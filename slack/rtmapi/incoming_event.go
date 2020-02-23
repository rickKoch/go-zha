package rtmapi

import (
	"encoding/json"
	"errors"
)

// Hello event is sent from slack when WebSocket connection is successfull.
type Hello struct {
	CommonEvent
}

// IncomingChannelEvent is any event in a channel
type IncomingChannelEvent struct {
	CommonEvent
	Channel string `json:"channel"`
}

// Message is message event on RTM
type Message struct {
	IncomingChannelEvent
	User      string    `json:"user"`
	Text      string    `json:"text"`
	TimeStamp TimeStamp `json:"ts"`
}

// TeamMigrationStarted is sent when chat group is migrated between servers.
type TeamMigrationStarted struct {
	CommonEvent
}

// Pong is given when client send Ping
type Pong struct {
	CommonEvent
	ReplyTo uint `json:"reply_to"`
}

// DecodedEvent is just an empty interface that marks decoded event.
type DecodedEvent interface{}

// DecodeEvent decodes given payload and converts this to corresponding event structure.
func DecodeEvent(input json.RawMessage) (DecodedEvent, error) {
	event := &CommonEvent{}
	if err := json.Unmarshal(input, event); err != nil {
		return nil, NewPayloadError(err.Error())
	}

	var mapping DecodedEvent

	switch event.Type {
	case HELLO:
		mapping = &Hello{}
	case MESSAGE:
		mapping = &Message{}
	case MIGRATION:
		mapping = &TeamMigrationStarted{}
	case PONG:
		mapping = &Pong{}
	case "":
		return nil, NewEventTypeError("type is not given" + string(input))
	default:
		return nil, NewUnknownEventTypeError("received unknwon event." + string(input))
	}

	if err := json.Unmarshal(input, mapping); err != nil {
		return nil, errors.New("error on JSON deserializing to mapped event" + string(input))
	}

	return mapping, nil
}
