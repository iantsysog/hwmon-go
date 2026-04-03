package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/iantsysog/hwmon-go/hwmon"
)

func SortReadingsStable(rs []hwmon.Reading) {
	sort.SliceStable(rs, func(i, j int) bool {
		if rs[i].Kind != rs[j].Kind {
			return rs[i].Kind < rs[j].Kind
		}
		if rs[i].Name != rs[j].Name {
			return rs[i].Name < rs[j].Name
		}
		if rs[i].KeyOrID != rs[j].KeyOrID {
			return rs[i].KeyOrID < rs[j].KeyOrID
		}
		return rs[i].Source < rs[j].Source
	})
}

func WriteReadingsJSON(w io.Writer, rs []hwmon.Reading) error {
	if w == nil {
		return fmt.Errorf("nil writer")
	}
	SortReadingsStable(rs)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rs)
}

func PrintWarning(warn error) {
	if warn != nil {
		_, _ = fmt.Fprintln(os.Stderr, "warning:", warn)
	}
}
