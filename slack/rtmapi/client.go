package rtmapi

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

// Client struct
type Client struct {
	PaylaodDecoder func(json.RawMessage) (DecodedEvent, error)
}

// NewClient creates new Client
func NewClient() *Client {
	return &Client{PaylaodDecoder: DefaultPayloadDecoder}
}

// Connect web socket connection
func (c *Client) Connect(url string) (*websocket.Conn, error) {
	return websocket.Dial(url, "", "http://localhost")
}

// DefaultPayloadDecoder decodes given paylaods, which includes various kinds of events.
func DefaultPayloadDecoder(payload json.RawMessage) (DecodedEvent, error) {
	decodedEvent, eventDecodedError := DecodeEvent(payload)

	if _, ok := eventDecodedError.(*EventTypeError); ok {
		// Check the reply status
		reply, err := DecodeReply(payload)
		if err != nil {
			return nil, NewPayloadError(err.Error())
		}

		if !*reply.OK {
			return nil, NewReplyStatusError(reply)
		}

		return nil, nil
	}

	if eventDecodedError != nil {
		return nil, NewPayloadError(eventDecodedError.Error())
	}

	return decodedEvent, nil
}
