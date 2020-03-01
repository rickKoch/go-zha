package redis

import "go.uber.org/zap"

// Option test
type Option func(*Config) error

// WithConfig test
func WithConfig(newConf Config) Option {
	return func(oldConf *Config) error {
		oldConf.Addr = newConf.Addr
		oldConf.Key = newConf.Key
		oldConf.Password = newConf.Password
		oldConf.DB = newConf.DB
		oldConf.Logger = newConf.Logger
		return nil
	}
}

// WithLogger test
func WithLogger(logger *zap.Logger) Option {
	return func(conf *Config) error {
		conf.Logger = logger
		return nil
	}
}

// WithKey test
func WithKey(key string) Option {
	return func(conf *Config) error {
		conf.Key = key
		return nil
	}
}
