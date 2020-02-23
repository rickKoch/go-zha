package rtmapi

// EventType is the type of the event sent from slack.
type EventType string

const (
	// UNKNOWN event type
	UNKNOWN = "unknown"
	// HELLO event type
	HELLO = "hello"
	// MESSAGE event type
	MESSAGE = "message"
	// MIGRATION is team_migration_started event type
	MIGRATION = "team_migration_started"
	// PING event type
	PING = "ping"
	// PONG event type
	PONG = "pong"
)

// CommonEvent have common fields on incoming/outgoing events
type CommonEvent struct {
	Type EventType `json:"type,omitempty"`
}
