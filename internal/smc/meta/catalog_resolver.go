package meta

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
)

type catalogResolver struct {
	cat *catalog.Catalog
}

func CatalogResolver(cat *catalog.Catalog) Resolver { return catalogResolver{cat: cat} }

func (r catalogResolver) lookup(key string) (Meta, bool) {
	if r.cat == nil {
		return Meta{}, false
	}
	e, ok := r.cat.Lookup(key)
	if !ok {
		return Meta{}, false
	}

	name := e.Name
	if name == "" {
		name = key
	}

	kind := hwmon.KindOther
	if e.Kind != "" {
		kind = hwmon.Kind(e.Kind)
	}
	return Meta{
		Name: name,
		Kind: kind,
		Unit: e.Unit,
		Hint: e.DataTypeHint,
	}, true
}

func (r catalogResolver) Lookup(ctx context.Context, key string) (Meta, bool) {
	if err := ctx.Err(); err != nil {
		return Meta{}, false
	}
	return r.lookup(key)
}
