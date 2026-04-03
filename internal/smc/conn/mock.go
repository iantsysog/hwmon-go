package conn

import (
	"sync"

	"github.com/iantsysog/hwmon-go/internal/smc/model"
)

type MockConnection struct {
	data map[string]model.SMCVal
	mu   sync.RWMutex
}

func NewMockConnection() Connection {
	return &MockConnection{
		data: make(map[string]model.SMCVal),
	}
}

func (c *MockConnection) WriteVal(key, dataType string, value []byte) error {
	if len(key) != 4 {
		return model.ErrInvalidKey
	}
	in := make([]byte, len(value))
	copy(in, value)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = model.SMCVal{Key: key, DataType: dataType, Bytes: in}
	return nil
}

func (c *MockConnection) Open() error {
	return nil
}

func (c *MockConnection) Close() error {
	return nil
}

func (c *MockConnection) Read(key string) (model.SMCVal, error) {
	if len(key) != 4 {
		return model.SMCVal{}, model.ErrInvalidKey
	}

	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return model.SMCVal{}, model.ErrNoDataForKey
	}

	out := make([]byte, len(v.Bytes))
	copy(out, v.Bytes)

	return model.SMCVal{
		Key:      key,
		DataType: v.DataType,
		Bytes:    out,
	}, nil
}

func (c *MockConnection) Write(key string, value []byte) error {
	if len(key) != 4 {
		return model.ErrInvalidKey
	}

	in := make([]byte, len(value))
	copy(in, value)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = model.SMCVal{Key: key, DataType: "hex_", Bytes: in}

	return nil
}
