package codec

import (
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestDecodeIntegersBigEndian(t *testing.T) {
	v, err := Decode("ui16", []byte{0x12, 0x34})
	test.NoError(t, err)
	test.Eq(t, uint16(0x1234), v.(uint16))

	v, err = Decode("si16", []byte{0xff, 0xfe})
	test.NoError(t, err)
	test.Eq(t, int16(-2), v.(int16))
}

func TestEncodeIntegersBigEndian(t *testing.T) {
	raw, err := Encode("ui32", uint32(0x01020304))
	test.NoError(t, err)
	test.Eq(t, []byte{0x01, 0x02, 0x03, 0x04}, raw)

	raw, err = Encode("si8 ", int8(-1))
	test.NoError(t, err)
	test.Eq(t, []byte{0xff}, raw)
}

func TestDecodeFlag(t *testing.T) {
	v, err := Decode("flag", []byte{0x00})
	test.NoError(t, err)
	test.Eq(t, false, v)

	v, err = Decode("flag", []byte{0x01})
	test.NoError(t, err)
	test.Eq(t, true, v)
}

func TestDecodeFixedPoint(t *testing.T) {
	raw := make([]byte, 2)
	binary.BigEndian.PutUint16(raw, uint16(int16(25*256)))
	v, err := Decode("sp78", raw)
	test.NoError(t, err)
	test.Eq(t, float64(25), v.(float64))

	binary.BigEndian.PutUint16(raw, 100)
	v, err = Decode("fpe2", raw)
	test.NoError(t, err)
	test.Eq(t, float64(25), v.(float64))
}

func TestEncodeFixedPointRoundTrip(t *testing.T) {
	raw, err := Encode("sp78", 12.5)
	test.NoError(t, err)
	v, err := Decode("sp78", raw)
	test.NoError(t, err)
	test.Eq(t, 12.5, v)
}

func TestFloat32(t *testing.T) {
	raw, err := Encode("flt ", float32(1.5))
	test.NoError(t, err)
	v, err := Decode("flt ", raw)
	test.NoError(t, err)
	test.Eq(t, float32(1.5), v.(float32))

	raw, err = Encode("flt ", float32(math.NaN()))
	test.NoError(t, err)
	v, err = Decode("flt ", raw)
	test.NoError(t, err)
	test.True(t, math.IsNaN(float64(v.(float32))))
}

func TestFixedPointVariants(t *testing.T) {
	raw, err := Encode("fp4c", 1.0)
	test.NoError(t, err)
	v, err := Decode("fp4c", raw)
	test.NoError(t, err)
	test.Eq(t, float64(1.0), v.(float64))

	raw, err = Encode("sp5a", -1.0)
	test.NoError(t, err)
	v, err = Decode("sp5a", raw)
	test.NoError(t, err)
	test.Eq(t, float64(-1.0), v.(float64))
}

func TestDecodeStringTypes(t *testing.T) {
	v, err := Decode("ch8*", []byte("abc"))
	test.NoError(t, err)
	test.Eq(t, "abc", v)

	v, err = Decode("cstr", []byte{'a', 0, 'b'})
	test.NoError(t, err)
	test.Eq(t, "a", v)
}

func TestInvalidLength(t *testing.T) {
	_, err := Decode("ui16", []byte{0x01})
	test.Error(t, err)
	var e ErrInvalidLength
	test.True(t, errors.As(err, &e))
	test.Eq(t, "ui16", e.DataType)
}
