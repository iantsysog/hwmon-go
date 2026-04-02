package smc

import (
	"sync"
)

type MockConnection struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMockConnection() Connection {
	return &MockConnection{
		data: make(map[string][]byte),
	}
}

func (c *MockConnection) Open() error {
	return nil
}

func (c *MockConnection) Close() error {
	return nil
}

func (c *MockConnection) Read(key string) (SMCVal, error) {
	if len(key) != 4 {
		return SMCVal{}, ErrInvalidKey
	}

	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return SMCVal{}, ErrNoDataForKey
	}

	out := make([]byte, len(v))
	copy(out, v)

	return SMCVal{
		Key:      key,
		DataType: "hex_",
		Bytes:    out,
	}, nil
}

func (c *MockConnection) Write(key string, value []byte) error {
	if len(key) != 4 {
		return ErrInvalidKey
	}

	in := make([]byte, len(value))
	copy(in, value)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = in

	return nil
}
