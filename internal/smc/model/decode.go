package model

import (
	"slices"

	"github.com/iantsysog/hwmon-go/internal/smc/codec"
)

func (v SMCVal) Decode() (codec.Value, error) {
	decoded, err := codec.Decode(v.DataType, v.Bytes)
	if err != nil {
		return codec.Value{}, err
	}
	return codec.Value{
		Key:      v.Key,
		DataType: v.DataType,
		Raw:      slices.Clone(v.Bytes),
		Decoded:  decoded,
	}, nil
}
