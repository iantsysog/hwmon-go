package codec

import (
	"encoding/binary"
	"math"
)

type fixed16Spec struct {
	dataType string
	div      float64
	signed   bool
}

func registerAppleFixed16() {
	specs := []fixed16Spec{
		{dataType: "fp1f", div: 32768.0, signed: false},
		{dataType: "fp2e", div: 16384.0, signed: false},
		{dataType: "fp3d", div: 8192.0, signed: false},
		{dataType: "fp4c", div: 4096.0, signed: false},
		{dataType: "fp5b", div: 2048.0, signed: false},
		{dataType: "fp6a", div: 1024.0, signed: false},
		{dataType: "fp79", div: 512.0, signed: false},
		{dataType: "fp88", div: 256.0, signed: false},
		{dataType: "fpa6", div: 64.0, signed: false},
		{dataType: "fpc4", div: 16.0, signed: false},
		{dataType: "fpe2", div: 4.0, signed: false},

		{dataType: "sp1e", div: 16384.0, signed: true},
		{dataType: "sp2d", div: 8192.0, signed: true},
		{dataType: "sp3c", div: 4096.0, signed: true},
		{dataType: "sp4b", div: 2048.0, signed: true},
		{dataType: "sp5a", div: 1024.0, signed: true},
		{dataType: "sp69", div: 512.0, signed: true},
		{dataType: "sp78", div: 256.0, signed: true},
		{dataType: "sp87", div: 128.0, signed: true},
		{dataType: "sp96", div: 64.0, signed: true},
		{dataType: "spa5", div: 32.0, signed: true},
		{dataType: "spb4", div: 16.0, signed: true},
		{dataType: "spf0", div: 1.0, signed: true},
	}

	for _, s := range specs {
		Register(s.dataType, fixed16Decoder(s.dataType, s.div, s.signed), fixed16Encoder(s.dataType, s.div, s.signed))
	}

	Register("ioft", func(raw []byte) (any, error) {
		if len(raw) != 8 {
			return nil, ErrInvalidLength{DataType: "ioft", Got: len(raw), Want: 8}
		}
		u := binary.LittleEndian.Uint64(raw)
		return float64(u) / 65536.0, nil
	}, func(v any) ([]byte, error) {
		f, ok := toFloat(v)
		if !ok {
			return nil, ErrUnsupportedValueType{DataType: "ioft", Value: v}
		}
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, ErrOutOfRange{DataType: "ioft", Value: v}
		}

		scaled := math.Round(f * 65536.0)
		if scaled < 0 || scaled > float64(^uint64(0)) {
			return nil, ErrOutOfRange{DataType: "ioft", Value: v}
		}
		out := make([]byte, 8)
		binary.LittleEndian.PutUint64(out, uint64(scaled))
		return out, nil
	})
}

func fixed16Decoder(dataType string, div float64, signed bool) Decoder {
	return func(raw []byte) (any, error) {
		if len(raw) != 2 {
			return nil, ErrInvalidLength{DataType: dataType, Got: len(raw), Want: 2}
		}
		u := binary.BigEndian.Uint16(raw)
		if signed {
			return float64(int16(u)) / div, nil
		}
		return float64(u) / div, nil
	}
}

func fixed16Encoder(dataType string, div float64, signed bool) Encoder {
	return func(v any) ([]byte, error) {
		f, ok := toFloat(v)
		if !ok {
			return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
		}
		scaled := math.Round(f * div)
		out := make([]byte, 2)
		if signed {
			if scaled < math.MinInt16 || scaled > math.MaxInt16 {
				return nil, ErrOutOfRange{DataType: dataType, Value: v}
			}
			binary.BigEndian.PutUint16(out, uint16(int16(scaled)))
			return out, nil
		}
		if scaled < 0 || scaled > math.MaxUint16 {
			return nil, ErrOutOfRange{DataType: dataType, Value: v}
		}
		binary.BigEndian.PutUint16(out, uint16(scaled))
		return out, nil
	}
}
