//go:build smc

package smc

import (
	"github.com/iantsysog/hwmon-go/internal/smc/catalog"
	"github.com/iantsysog/hwmon-go/internal/smc/conn"
)

func withConnection(cn conn.Connection) Option { return func(c *config) { c.conn = cn } }

func withCatalog(cat *catalog.Catalog) Option { return func(c *config) { c.cat = cat } }
