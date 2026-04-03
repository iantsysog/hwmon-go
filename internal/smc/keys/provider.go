package keys

import (
	"context"
	"slices"
	"sort"
)

type Provider interface {
	Keys(ctx context.Context) ([]string, error)
}

type staticProvider struct {
	keys []string
}

func Static(keys []string) Provider {
	return staticProvider{keys: slices.Clone(keys)}
}

func (p staticProvider) Keys(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return slices.Clone(p.keys), nil
}

type catalogKeys interface {
	Keys() []string
}

type catalogProvider struct {
	cat catalogKeys
}

func FromCatalog(cat catalogKeys) Provider {
	return catalogProvider{cat: cat}
}

func (p catalogProvider) Keys(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p.cat == nil {
		return nil, nil
	}
	return slices.Clone(p.cat.Keys()), nil
}

type concatProvider struct {
	ps []Provider
}

func Concat(ps ...Provider) Provider {
	cp := append([]Provider(nil), ps...)
	return concatProvider{ps: cp}
}

func (p concatProvider) Keys(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, 64)
	var out []string
	for _, src := range p.ps {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if src == nil {
			continue
		}
		ks, err := src.Keys(ctx)
		if err != nil {
			return nil, err
		}
		for _, k := range ks {
			if k == "" {
				continue
			}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out, nil
}

type filterProvider struct {
	p  Provider
	fn func(string) bool
}

func Filter(p Provider, fn func(string) bool) Provider {
	return filterProvider{p: p, fn: fn}
}

func (p filterProvider) Keys(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p.p == nil {
		return nil, nil
	}
	ks, err := p.p.Keys(ctx)
	if err != nil {
		return nil, err
	}
	if p.fn == nil {
		return slices.Clone(ks), nil
	}
	out := make([]string, 0, len(ks))
	for _, k := range ks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if p.fn(k) {
			out = append(out, k)
		}
	}
	return out, nil
}
