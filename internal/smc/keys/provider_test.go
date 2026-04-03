package keys

import (
	"context"
	"slices"
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

type fakeCat struct {
	ks []string
}

func (f fakeCat) Keys() []string { return slices.Clone(f.ks) }

func TestStatic(t *testing.T) {
	p := Static([]string{"B", "A"})
	ks, err := p.Keys(context.Background())
	test.NoError(t, err)
	test.Eq(t, []string{"B", "A"}, ks)
}

func TestFromCatalogNil(t *testing.T) {
	p := FromCatalog(nil)
	ks, err := p.Keys(context.Background())
	test.NoError(t, err)
	test.Eq(t, []string(nil), ks)
}

func TestConcatDedupSort(t *testing.T) {
	p := Concat(
		Static([]string{"B", "A"}),
		FromCatalog(fakeCat{ks: []string{"A", "C"}}),
	)
	ks, err := p.Keys(context.Background())
	test.NoError(t, err)
	test.Eq(t, []string{"A", "B", "C"}, ks)
}

func TestFilter(t *testing.T) {
	p := Filter(Static([]string{"TC0C", "FNum"}), func(k string) bool { return k == "FNum" })
	ks, err := p.Keys(context.Background())
	test.NoError(t, err)
	test.Eq(t, []string{"FNum"}, ks)
}
