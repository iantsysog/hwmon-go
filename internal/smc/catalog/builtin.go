//go:build catalog

package catalog

import (
	"bytes"
	"embed"
	"sync"
)

//go:embed generated.json
var builtinFS embed.FS

var (
	builtinOnce sync.Once
	builtinCat  *Catalog
)

func builtin() *Catalog {
	builtinOnce.Do(func() {
		b, err := builtinFS.ReadFile("generated.json")
		if err != nil {
			return
		}
		c, err := Load(bytes.NewReader(b))
		if err != nil {
			return
		}
		builtinCat = c
	})
	return builtinCat
}
