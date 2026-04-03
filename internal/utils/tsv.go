package utils

import (
	"fmt"
	"io"
)

func WriteTSVRow(w io.Writer, cols ...string) {
	if w == nil {
		return
	}
	for i, c := range cols {
		if i > 0 {
			_, _ = w.Write([]byte{'\t'})
		}
		_, _ = io.WriteString(w, c)
	}
	_, _ = fmt.Fprintln(w)
}
