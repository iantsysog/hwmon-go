package meta

import (
	"context"
	"strings"
	"testing"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestFallbackResolver(t *testing.T) {
	r := FallbackResolver()
	m, ok := r.Lookup(context.Background(), "TC0C")
	test.True(t, ok)
	test.Eq(t, "TC0C", m.Name)
	test.Eq(t, hwmon.KindOther, m.Kind)
}

func TestCatalogResolver(t *testing.T) {
	in := `{"version":1,"entries":[{"key":"TC0C","name":"CPU","kind":"temp","unit":"°C","dataTypeHint":"sp78"}]}`
	cat, err := catalog.Load(strings.NewReader(in))
	test.NoError(t, err)

	r := CatalogResolver(cat)
	m, ok := r.Lookup(context.Background(), "TC0C")
	test.True(t, ok)
	test.Eq(t, "CPU", m.Name)
	test.Eq(t, hwmon.KindTemp, m.Kind)
	test.Eq(t, "°C", m.Unit)
	test.Eq(t, "sp78", m.Hint)
}

func TestChain(t *testing.T) {
	in := `{"version":1,"entries":[{"key":"TC0C","name":"CPU","kind":"temp","unit":"°C","dataTypeHint":"sp78"}]}`
	cat, err := catalog.Load(strings.NewReader(in))
	test.NoError(t, err)

	r := Chain(CatalogResolver(cat), FallbackResolver())
	m, ok := r.Lookup(context.Background(), "UNKNOWN")
	test.True(t, ok)
	test.Eq(t, "UNKNOWN", m.Name)
}
