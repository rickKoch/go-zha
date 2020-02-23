package rtmapi

import "fmt"

// PayloadError represents an error that given JSON payload is not properly formatted.
type PayloadError struct {
	Err string
}

// Error returns the error string
func (e *PayloadError) Error() string {
	return e.Err;
}

// NewPayloadError creates new PayloadError
func NewPayloadError(str string) *PayloadError {
	return &PayloadError{Err: str}
}

// EventTypeError is error on JSON string parsing.
type EventTypeError struct {
	Err string
}

// Error returns the error string
func (e *EventTypeError) Error() string {
	return e.Err;
}

// NewEventTypeError create new EventTypeError
func NewEventTypeError(e string) *EventTypeError {
	return &EventTypeError{Err:e}
}

// UnknownEventTypeError is returned when given event's type is undefined
type UnknownEventTypeError struct {
	Err string
}

// Error returns the error string
func (e *UnknownEventTypeError) Error() string {
	return e.Err;
}

// NewUnknownEventTypeError creates new UnknownEventTypeError
func NewUnknownEventTypeError(e string) *UnknownEventTypeError {
	return &UnknownEventTypeError{Err: e}
}

// ReplyStatusError is returned when given WebSocketReply payload gives status error
type ReplyStatusError struct {
	Reply *WebSocketReply
}

// Error returns its error string
func (e *ReplyStatusError) Error() string {
	return fmt.Sprintf("error on previous message posting %#v", e.Reply)
}

// NewReplyStatusError create new ReplyStatusError
func NewReplyStatusError(reply *WebSocketReply) *ReplyStatusError {
	return &ReplyStatusError{Reply: reply}
}
