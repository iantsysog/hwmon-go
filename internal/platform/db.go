//go:build platform

package platform

import (
	"embed"
	"encoding/json"
	"sync"
)

//go:embed products.json
var productsFS embed.FS

type product struct {
	Family string `json:"family"`
}

var families = sync.OnceValue(func() map[string]string {
	b, err := productsFS.ReadFile("products.json")
	if err != nil {
		return map[string]string{}
	}

	var parsed map[string]product
	if err := json.Unmarshal(b, &parsed); err != nil {
		return map[string]string{}
	}

	out := make(map[string]string, len(parsed))
	for model, p := range parsed {
		if p.Family == "" {
			continue
		}
		out[model] = p.Family
	}
	return out
})

func lookupFamily(model string) (string, bool) {
	f, ok := families()[model]
	if !ok || f == "" {
		return "", false
	}
	return f, true
}
