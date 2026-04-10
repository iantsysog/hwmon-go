package smc

import "slices"

type config struct {
	keys          []string
	useCatalog    bool
	emitUndecoded bool
}

type Option func(*config)

func WithKeys(keys ...string) Option {
	return func(c *config) {
		c.keys = slices.Clone(keys)
	}
}

func WithCatalogEnabled(enabled bool) Option { return func(c *config) { c.useCatalog = enabled } }

func WithEmitUndecoded(enabled bool) Option {
	return func(c *config) { c.emitUndecoded = enabled }
}
