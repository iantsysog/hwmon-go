package meta

import (
	"context"

	"github.com/iantsysog/hwmon-go/hwmon"
)

type Meta struct {
	Name string
	Kind hwmon.Kind
	Unit string
	Hint string
}

type Resolver interface {
	Lookup(ctx context.Context, key string) (Meta, bool)
}

type fallbackResolver struct{}

func FallbackResolver() Resolver { return fallbackResolver{} }

type chainResolver struct {
	primary  Resolver
	fallback Resolver
}

func Chain(primary, fallback Resolver) Resolver {
	return chainResolver{primary: primary, fallback: fallback}
}

func (r chainResolver) Lookup(ctx context.Context, key string) (Meta, bool) {
	if err := ctx.Err(); err != nil {
		return Meta{}, false
	}
	if r.primary != nil {
		if m, ok := r.primary.Lookup(ctx, key); ok {
			return m, true
		}
	}
	if r.fallback != nil {
		if m, ok := r.fallback.Lookup(ctx, key); ok {
			return m, true
		}
	}
	return Meta{}, false
}

func StaticFallback(key string) Meta {
	return Meta{
		Name: key,
		Kind: hwmon.KindOther,
		Unit: "",
		Hint: "",
	}
}

func (fallbackResolver) Lookup(ctx context.Context, key string) (Meta, bool) {
	if err := ctx.Err(); err != nil {
		return Meta{}, false
	}
	return StaticFallback(key), true
}
