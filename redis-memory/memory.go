package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gitlab.com/kochevRisto/go-zha"
	"go.uber.org/zap"
)

// Config test
type Config struct {
	Addr     string
	Key      string
	Password string
	DB       int
	Logger   *zap.Logger
}

type memory struct {
	logger *zap.Logger
	Client *redis.Client
	hkey   string
}

// Memory test
func Memory(addr string, opts ...Option) zha.Option {
	return func(b *zha.Bot) error {
		conf := Config{Addr: addr}
		for _, opt := range opts {
			err := opt(&conf)
			if err != nil {
				return err
			}
		}

		if b.Logger != nil {
			opts = append(opts, WithLogger(b.Logger.Named("memory")))
		}

		memory, err := NewMemory(conf)
		if err != nil {
			return err
		}

		b.Memory = memory
		return nil
	}
}

// NewMemory test
func NewMemory(conf Config) (zha.Memory, error) {
	if conf.Logger == nil {
		conf.Logger = zap.NewNop()
	}

	if conf.Key == "" {
		conf.Key = "risto-bot"
	}

	memory := &memory{
		logger: conf.Logger,
		hkey:   conf.Key,
	}

	memory.logger.Debug("Connecting to redis memory",
		zap.String("addr", conf.Addr),
		zap.String("key", memory.hkey),
	)

	memory.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := memory.Client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping redis")
	}

	memory.logger.Info("Memory initialized successfully")
	return memory, nil
}

// Set test
func (m *memory) Set(key, value string) error {
	m.logger.Debug("Writing data to memory", zap.String("key", key))
	resp := m.Client.HSet(m.hkey, key, value)
	return resp.Err()
}

// Get test
func (m *memory) Get(key string) (string, bool, error) {
	m.logger.Debug("Retriving data from memory", zap.String("key", key))
	res, err := m.Client.HGet(m.hkey, key).Result()
	switch {
	case err == redis.Nil:
		return "", false, nil
	case err != nil:
		return "", false, err
	default:
		return res, true, nil
	}
}

// Delete test
func (m *memory) Delete(key string) (bool, error) {
	m.logger.Debug("Deleting data from memory", zap.String("key", key))
	res, err := m.Client.HDel(m.hkey, key).Result()
	return res > 0, err
}

// Memories test
func (m *memory) Memories() (map[string]string, error) {
	return m.Client.HGetAll(m.hkey).Result()
}

// Close test
func (m *memory) Close() error {
	return m.Client.Close()
}
