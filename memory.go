package zha

import "sync"

// Memory interface
type Memory interface {
	Set(key, value string) error
	Get(key string) (string, bool, error)
	Delete(key string) (bool, error)
	Memories() (map[string]string, error)
	Close() error
}

// InMemory struct
type InMemory struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewInMemory returns new InMemory memory
func NewInMemory() *InMemory {
	return &InMemory{
		data: map[string]string{},
	}
}

// Set sets value to the memory
func (m *InMemory) Set(key, value string) error {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()

	return nil
}

// Get gets data from the memory
func (m *InMemory) Get(key string) (string, bool, error) {
	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()

	return value, ok, nil
}

// Delete removes data from the memory
func (m *InMemory) Delete(key string) (bool, error) {
	m.mu.Lock()
	_, ok := m.data[key]
	delete(m.data, key)
	m.mu.Unlock()

	return ok, nil
}

// Close should close connection to memroy
func (m *InMemory) Close() error {
	m.mu.Lock()
	m.data = map[string]string{}
	m.mu.Unlock()

	return nil
}

// Memories return memory data
func (m *InMemory) Memories() (map[string]string, error) {
	m.mu.RLock()
	memories := make(map[string]string, len(m.data))
	for k, v := range m.data {
		memories[k] = v
	}
	m.mu.RUnlock()

	return memories, nil
}
