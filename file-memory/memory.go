package file

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/pkg/errors"
	"gitlab.com/kochevRisto/go-zha"
	"go.uber.org/zap"
)

// Memory is memory struct
type Memory struct {
	path   string
	logger *zap.Logger
	mu     sync.RWMutex
	data   map[string]string
}

// NewMemory creates new file memory
func NewMemory(path string, opts ...Option) (*Memory, error) {
	memory := &Memory{
		path: path,
		data: map[string]string{},
	}

	for _, opt := range opts {
		err := opt(memory)
		if err != nil {
			return nil, err
		}

	}

	if memory.logger == nil {
		memory.logger = zap.NewNop()
	}

	memory.logger.Debug("Opening memor file", zap.String("path", path))

	f, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		memory.logger.Debug("File does not exist. Continue with empty memory", zap.String("path", path))
	case err != nil:
		return nil, errors.Wrap(err, "failed to open file")
	default:
		memory.logger.Debug("Decoding JSON from memory file", zap.String("path", path))
		err := json.NewDecoder(f).Decode(&memory.data)
		_ = f.Close()
		if err != nil {
			return nil, errors.Wrap(err, "failed decode data as JSON")
		}
	}

	memory.logger.Info(
		"Memory init successfully",
		zap.String("path", path),
		zap.Int("num_memories", len(memory.data)),
	)

	return memory, nil
}

// MemoryOption sets the memory options
func MemoryOption(path string) zha.Option {
	return func(b *zha.Bot) error {
		var opts []Option
		if b.Logger != nil {
			opts = append(opts, WithLogger(b.Logger.Named("memory")))
		}

		memory, err := NewMemory(path, opts...)
		if err != nil {
			return err
		}

		b.Memory = memory
		return nil
	}
}

// Set test
func (m *Memory) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		return errors.New("memory was alredy shut down")
	}

	m.logger.Debug("Writing data to memory", zap.String("key", key))
	m.data[key] = value
	err := m.persist()

	return err
}

// Get test
func (m *Memory) Get(key string) (string, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data == nil {
		return "", false, errors.New("memory was already shut down")
	}

	m.logger.Debug("Retrieving data from memory", zap.String("key", key))
	value, ok := m.data[key]
	return value, ok, nil
}

// Delete test
func (m *Memory) Delete(key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		return false, errors.New("memory was already shut down")
	}

	m.logger.Debug("Deleting data from memory", zap.String("key", key))
	_, ok := m.data[key]
	delete(m.data, key)
	err := m.persist()

	return ok, err
}

// Memories test
func (m *Memory) Memories() (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data == nil {
		return nil, errors.New("memory was already shut down")
	}

	mem := make(map[string]string, len(m.data))
	for k, v := range m.data {
		mem[k] = v
	}

	return mem, nil
}

// Close test
func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		return errors.New("memory was already shut down")
	}

	m.logger.Debug("shuting down memory")
	m.data = nil

	return nil
}

func (m *Memory) persist() error {
	f, err := os.OpenFile(m.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return errors.Wrap(err, "failed to open file to persist data")
	}

	err = json.NewEncoder(f).Encode(m.data)
	if err != nil {
		_ = f.Close()
		return errors.Wrap(err, "failed to encode data as JSON")
	}

	err = f.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close file; data might not have been fully persisted to disk")
	}

	return nil
}
