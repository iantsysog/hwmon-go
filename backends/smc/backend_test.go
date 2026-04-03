//go:build smc

package smc

import (
	"context"
	"strings"
	"testing"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
	"github.com/iantsysog/hwmon-go/internal/smc/conn"
	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestBackend_Collect_UsesCatalogKeysAndMetadata(t *testing.T) {
	c := conn.NewMockConnection().(*conn.MockConnection)
	test.NoError(t, c.WriteVal("TC0C", "sp78", []byte{0x19, 0x00}))

	in := `{"version":1,"entries":[{"key":"TC0C","name":"CPU Core Temp","kind":"temp","unit":"°C","dataTypeHint":"sp78"}]}`
	cat, err := catalog.Load(strings.NewReader(in))
	test.NoError(t, err)

	b := New(
		withConnection(c),
		withCatalog(cat),
		WithCatalogEnabled(true),
	)

	var out []hwmon.Reading
	test.NoError(t, b.Collect(context.Background(), func(r hwmon.Reading) { out = append(out, r) }))
	test.Len(t, out, 1)
	test.Eq(t, hwmon.KindTemp, out[0].Kind)
	test.Eq(t, "CPU Core Temp", out[0].Name)
	test.Eq(t, "°C", out[0].Unit)
	test.Eq(t, "smc", out[0].Source)
	test.Eq(t, "TC0C", out[0].KeyOrID)
	test.Eq(t, "sp78", out[0].DataType)
	test.Eq(t, float64(25), out[0].Value)
	test.Len(t, out[0].Raw, 2)
}

func TestBackend_Collect_WithKeysNoCatalog(t *testing.T) {
	c := conn.NewMockConnection().(*conn.MockConnection)
	test.NoError(t, c.WriteVal("FNum", "ui8 ", []byte{0x02}))

	b := New(
		withConnection(c),
		WithCatalogEnabled(false),
		WithKeys("FNum"),
	)

	var out []hwmon.Reading
	test.NoError(t, b.Collect(context.Background(), func(r hwmon.Reading) { out = append(out, r) }))
	test.Len(t, out, 1)
	test.Eq(t, hwmon.KindOther, out[0].Kind)
	test.Eq(t, "FNum", out[0].Name)
	test.Eq(t, uint8(2), out[0].Value)
}
