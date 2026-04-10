package hwmon

import "sync"

type Registry interface {
	List() []Backend
}

type registry struct {
	backends []Backend
	mu       sync.RWMutex
}

func (r *registry) List() []Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Backend, len(r.backends))
	copy(out, r.backends)
	return out
}

func (r *registry) add(b Backend) {
	if b == nil {
		return
	}
	r.mu.Lock()
	r.backends = append(r.backends, b)
	r.mu.Unlock()
}

var defaultRegistry = &registry{}

func Register(b Backend) {
	defaultRegistry.add(b)
}

func DefaultRegistry() Registry { return defaultRegistry }
