package slack

import "go.uber.org/zap"

// Option is Slack options
type Option func(*Config) error

// WithLogger sets logger on the slack adapter
func WithLogger(logger *zap.Logger) Option {
	return func(conf *Config) error {
		conf.Logger = logger
		return nil
	}
}
