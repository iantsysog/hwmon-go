package model

import (
	"encoding/binary"
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestSMCValUintInt(t *testing.T) {
	v := SMCVal{Bytes: []byte{0x01, 0x00}}
	u, err := v.Uint(binary.LittleEndian)
	test.NoError(t, err)
	test.Eq(t, uint64(1), u)

	i, err := v.Int(binary.LittleEndian)
	test.NoError(t, err)
	test.Eq(t, int64(1), i)
}

func TestSMCValFloat32(t *testing.T) {
	v := SMCVal{Bytes: BytesFloat32(binary.LittleEndian, 12.5)}
	f, err := v.Float32LE()
	test.NoError(t, err)
	test.Eq(t, float32(12.5), f)
}

func TestSMCValCString(t *testing.T) {
	v := SMCVal{Bytes: []byte{'h', 'i', 0x00, 'x'}}
	test.Eq(t, "hi", v.CString())
	test.Eq(t, "hi\x00x", v.String())
}

func TestSMCValErrors(t *testing.T) {
	_, err := (SMCVal{}).Bool()
	test.Eq(t, ErrInsufficientBytes, err)

	_, err = (SMCVal{Bytes: []byte{0x01, 0x02, 0x03}}).Uint(binary.BigEndian)
	test.Eq(t, ErrUnsupportedSize, err)

	_, err = (SMCVal{Bytes: []byte{0x01, 0x02}}).Float32LE()
	test.Eq(t, ErrInsufficientBytes, err)
}
