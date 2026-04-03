//go:build hid

package hid

import (
	"context"
	"slices"
	"testing"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/hid/model"
	"github.com/iantsysog/hwmon-go/internal/hid/source"
	"github.com/iantsysog/hwmon-go/internal/test"
)

type fakeSource struct {
	rs []model.Reading
}

func (f fakeSource) Open() error  { return nil }
func (f fakeSource) Close() error { return nil }
func (f fakeSource) Readings() ([]model.Reading, error) {
	return slices.Clone(f.rs), nil
}

func TestBackend_Collect_MapsReadingsAndFiltersOther(t *testing.T) {
	var _ source.Source = fakeSource{}
	b := New(withSource(fakeSource{
		rs: []model.Reading{
			{Kind: model.KindOther, Name: "x", Value: 1, KeyOrID: "x"},
			{Kind: model.KindTemp, Name: "t", Unit: "°C", Value: 42, KeyOrID: "t"},
		},
	}))

	var out []hwmon.Reading
	test.NoError(t, b.Collect(context.Background(), func(r hwmon.Reading) { out = append(out, r) }))
	test.Len(t, out, 1)
	test.Eq(t, hwmon.KindTemp, out[0].Kind)
	test.Eq(t, "t", out[0].Name)
	test.Eq(t, float64(42), out[0].Value)
}

func TestBackend_Collect_EmitUnclassified(t *testing.T) {
	b := New(
		withSource(fakeSource{rs: []model.Reading{{Kind: model.KindOther, Name: "x", Value: 1, KeyOrID: "x"}}}),
		WithEmitUnclassified(true),
	)

	var out []hwmon.Reading
	test.NoError(t, b.Collect(context.Background(), func(r hwmon.Reading) { out = append(out, r) }))
	test.Len(t, out, 1)
}
