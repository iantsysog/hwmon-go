package hwmon

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type ErrorPolicy int

const (
	ErrorPolicyJoin ErrorPolicy = iota
	ErrorPolicyFailFast
)

type collectorConfig struct {
	registry Registry

	enabledNames  map[string]struct{}
	extraBackends []Backend

	filter func(Reading) bool
	less   func(a, b Reading) bool

	errorPolicy ErrorPolicy
}

type Option func(*collectorConfig)

func WithRegistry(r Registry) Option {
	return func(c *collectorConfig) {
		if r != nil {
			c.registry = r
		}
	}
}

func WithBackends(names ...string) Option {
	return func(c *collectorConfig) {
		if c.enabledNames == nil {
			c.enabledNames = make(map[string]struct{}, len(names))
		}
		for _, n := range names {
			if n == "" {
				continue
			}
			c.enabledNames[n] = struct{}{}
		}
	}
}

func WithBackend(b Backend) Option {
	return func(c *collectorConfig) {
		if b != nil {
			c.extraBackends = append(c.extraBackends, b)
		}
	}
}

func WithFilter(f func(Reading) bool) Option {
	return func(c *collectorConfig) { c.filter = f }
}

func WithSort(less func(a, b Reading) bool) Option {
	return func(c *collectorConfig) { c.less = less }
}

func WithErrorPolicy(p ErrorPolicy) Option {
	return func(c *collectorConfig) { c.errorPolicy = p }
}

func Collect(ctx context.Context, opts ...Option) ([]Reading, error) {
	cfg := collectorConfig{
		registry:    DefaultRegistry(),
		errorPolicy: ErrorPolicyJoin,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	var backends []Backend
	if cfg.registry != nil {
		backends = append(backends, cfg.registry.List()...)
	}
	backends = append(backends, cfg.extraBackends...)

	if len(cfg.enabledNames) > 0 {
		filtered := backends[:0]
		for _, b := range backends {
			if b == nil {
				continue
			}
			if _, ok := cfg.enabledNames[b.Name()]; ok {
				filtered = append(filtered, b)
			}
		}
		backends = filtered
	}

	if len(backends) == 0 {
		return nil, ErrNoBackends
	}

	var out []Reading
	var outMu sync.Mutex
	errs := make([]error, 0, len(backends))

	emit := func(r Reading) {
		if cfg.filter != nil && !cfg.filter(r) {
			return
		}
		outMu.Lock()
		out = append(out, r)
		outMu.Unlock()
	}

	for _, b := range backends {
		if b == nil {
			continue
		}
		name := b.Name()
		if err := b.Collect(ctx, emit); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			if cfg.errorPolicy == ErrorPolicyFailFast {
				break
			}
		}
	}

	less := cfg.less
	if less == nil {
		less = func(a, b Reading) bool {
			if a.Kind != b.Kind {
				return a.Kind < b.Kind
			}
			if a.Name != b.Name {
				return a.Name < b.Name
			}
			if a.KeyOrID != b.KeyOrID {
				return a.KeyOrID < b.KeyOrID
			}
			return a.Source < b.Source
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return less(out[i], out[j]) })

	if len(out) == 0 && len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return out, errors.Join(errs...)
}
