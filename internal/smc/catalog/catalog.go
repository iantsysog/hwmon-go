package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
)

type Entry struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Unit         string `json:"unit"`
	DataTypeHint string `json:"dataTypeHint"`
}

type fileFormat struct {
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
}

type Catalog struct {
	entries map[string]Entry
}

func (c *Catalog) Lookup(key string) (Entry, bool) {
	if c == nil || c.entries == nil {
		return Entry{}, false
	}
	e, ok := c.entries[key]
	return e, ok
}

func (c *Catalog) Keys() []string {
	if c == nil || len(c.entries) == 0 {
		return nil
	}
	out := make([]string, 0, len(c.entries))
	for k := range c.entries {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func Load(r io.Reader) (*Catalog, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}

	dec := json.NewDecoder(r)
	var ff fileFormat
	if err := dec.Decode(&ff); err != nil {
		return nil, fmt.Errorf("decode catalog: %w", err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, errors.New("decode catalog: trailing data")
		}
		return nil, fmt.Errorf("decode catalog: trailing data: %w", err)
	}
	if ff.Version != 1 {
		return nil, fmt.Errorf("unsupported catalog format version: %d", ff.Version)
	}
	c := &Catalog{entries: make(map[string]Entry, len(ff.Entries))}
	for _, e := range ff.Entries {
		if e.Key == "" {
			return nil, errors.New("catalog entry missing key")
		}
		if e.Name == "" {
			e.Name = e.Key
		}
		c.entries[e.Key] = e
	}
	return c, nil
}

func Builtin() *Catalog { return builtin() }
