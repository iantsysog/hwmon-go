//go:build !smc

package smc

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/model"
)

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
	_ = ctx
	_ = emit
	return model.ErrUnsupported
}
