//go:build !hid

package hid

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/hid/source"
)

type backend struct {
	cfg config
}

func New(opts ...Option) hwmon.Backend {
	cfg := config{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return &backend{cfg: cfg}
}

func (b *backend) Name() string { return "hid" }

func (b *backend) Collect(context.Context, func(hwmon.Reading)) error {
	return source.ErrUnsupported
}
