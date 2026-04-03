package smc

import (
	"slices"

	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
	"github.com/iantsysog/hwmon-go/internal/smc/conn"
)

type config struct {
	keys          []string
	useCatalog    bool
	emitUndecoded bool
	conn          conn.Connection
	cat           *catalog.Catalog
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
