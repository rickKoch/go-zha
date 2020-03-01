package file

import "go.uber.org/zap"

// Option helper type for memory
type Option func(memory *Memory) error

// WithLogger set logger options
func WithLogger(logger *zap.Logger) Option {
	return func(m *Memory) error {
		m.logger = logger
		return nil
	}
}
