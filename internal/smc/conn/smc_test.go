//go:build darwin && cgo

package conn

import (
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestSMC(t *testing.T) {
	c := New()
	if err := c.Open(); err != nil {
		t.Skipf("SMC open failed: %v", err)
	}
	defer func() { _ = c.Close() }()

	if err := c.Write("CH0B", []byte{0x0}); err != nil {
		t.Skipf("SMC write failed: %v", err)
	}

	v, err := c.Read("CH0B")
	test.NoError(t, err)
	test.Eq(t, "CH0B", v.Key)
	test.Eq(t, "hex_", v.DataType)
	test.Eq(t, []uint8{0x0}, v.Bytes)
}
