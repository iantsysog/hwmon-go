//go:build hid

package hid

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/hid/model"
	"github.com/iantsysog/hwmon-go/internal/hid/source"
)

func init() { hwmon.Register(New()) }

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

func (b *backend) Collect(ctx context.Context, emit func(hwmon.Reading)) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if emit == nil {
		return nil
	}

	src := b.cfg.src
	if src == nil {
		src = source.Default()
	}
	if err := src.Open(); err != nil {
		return err
	}
	defer src.Close()

	rs, err := src.Readings()
	if err != nil {
		return err
	}

	for _, r := range rs {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !b.cfg.emitUnclassified && r.Kind == model.KindOther {
			continue
		}
		emit(hwmon.Reading{
			Kind:     hwmon.Kind(r.Kind),
			Name:     r.Name,
			Unit:     r.Unit,
			Source:   "hid",
			KeyOrID:  r.KeyOrID,
			DataType: "",
			Raw:      nil,
			Value:    r.Value,
		})
	}
	return nil
}
