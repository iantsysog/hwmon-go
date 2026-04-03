//go:build smc

package smc

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
	"github.com/iantsysog/hwmon-go/internal/smc/collect"
	"github.com/iantsysog/hwmon-go/internal/smc/conn"
	"github.com/iantsysog/hwmon-go/internal/smc/keys"
	"github.com/iantsysog/hwmon-go/internal/smc/meta"
)

func init() { hwmon.Register(New()) }

type backend struct {
	cfg config
}

func New(opts ...Option) hwmon.Backend {
	cfg := config{useCatalog: true}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return &backend{cfg: cfg}
}

func (b *backend) Name() string { return "smc" }

func (b *backend) Collect(ctx context.Context, emit func(hwmon.Reading)) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if emit == nil {
		return nil
	}

	cn := b.cfg.conn
	if cn == nil {
		cn = conn.New()
	}

	kp := keys.Static(b.cfg.keys)
	mr := meta.FallbackResolver()
	if b.cfg.useCatalog {
		catUsed := b.cfg.cat
		if catUsed == nil {
			catUsed = catalog.Builtin()
		}
		if len(b.cfg.keys) == 0 {
			kp = keys.FromCatalog(catUsed)
		}
		mr = meta.Chain(meta.CatalogResolver(catUsed), meta.FallbackResolver())
	}

	return collect.Run(ctx, cn, kp, mr, emit, collect.Options{EmitUndecoded: b.cfg.emitUndecoded})
}
