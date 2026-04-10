package collect

import (
	"context"
	"sort"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/codec"
	"github.com/iantsysog/hwmon-go/internal/smc/conn"
	"github.com/iantsysog/hwmon-go/internal/smc/keys"
	"github.com/iantsysog/hwmon-go/internal/smc/meta"
	"github.com/iantsysog/hwmon-go/internal/smc/model"
)

type Options struct {
	EmitUndecoded bool
}

func Run(ctx context.Context, c conn.Connection, kp keys.Provider, mr meta.Resolver, emit func(hwmon.Reading), opts Options) error {
	if c == nil {
		return model.ErrUnsupported
	}
	if kp == nil {
		return nil
	}
	if emit == nil {
		return nil
	}

	ks, err := kp.Keys(ctx)
	if err != nil {
		return err
	}
	if len(ks) == 0 {
		return nil
	}
	sort.Strings(ks)

	if err := c.Open(); err != nil {
		return err
	}
	defer func() { _ = c.Close() }()

	for _, key := range ks {
		if err := ctx.Err(); err != nil {
			return err
		}
		v, err := c.Read(key)
		if err != nil {
			continue
		}

		m := meta.StaticFallback(key)
		if mr != nil {
			if mm, ok := mr.Lookup(ctx, key); ok {
				m = mm
			}
		}

		var raw []byte
		if len(v.Bytes) > 0 {
			raw = make([]byte, len(v.Bytes))
			copy(raw, v.Bytes)
		}

		dataTypeOut := v.DataType
		if m.Hint != "" {
			dataTypeOut = m.Hint
		}

		var decoded any
		decoded, decErr := codec.Decode(dataTypeOut, raw)
		if decErr != nil && dataTypeOut != v.DataType {
			decoded, decErr = codec.Decode(v.DataType, raw)
			dataTypeOut = v.DataType
		}
		if decErr != nil {
			if !opts.EmitUndecoded {
				continue
			}
			decoded = nil
		}

		emit(hwmon.Reading{
			Kind:     m.Kind,
			Name:     m.Name,
			Unit:     m.Unit,
			Source:   "smc",
			KeyOrID:  key,
			DataType: dataTypeOut,
			Raw:      raw,
			Value:    decoded,
		})
	}

	return nil
}
