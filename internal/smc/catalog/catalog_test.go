package catalog

import (
	"strings"
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestLoadV1(t *testing.T) {
	in := `{"version":1,"entries":[{"key":"TC0C","name":"CPU","kind":"temp","unit":"°C","dataTypeHint":"sp78"}]}`
	c, err := Load(strings.NewReader(in))
	test.NoError(t, err)

	e, ok := c.Lookup("TC0C")
	test.True(t, ok)
	test.Eq(t, "CPU", e.Name)
	test.Eq(t, "temp", e.Kind)
	test.Eq(t, "°C", e.Unit)
	test.Eq(t, "sp78", e.DataTypeHint)
}
