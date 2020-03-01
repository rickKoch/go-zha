package zha

// Option type
type Option func(*Bot) error

// WithMemory sets a memory on a bot
func WithMemory(memory Memory) Option {
	return func(b *Bot) error {
		b.Memory = memory
		return nil
	}
}
