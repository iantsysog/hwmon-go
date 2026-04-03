//go:build darwin && cgo && eventsystem

package source

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework IOKit
#include <stdlib.h>
#include "eventsystem.h"
*/
import "C"

import (
	"bufio"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/iantsysog/hwmon-go/internal/hid/model"
)

type eventSystemSource struct{}

func defaultEventSystemSource() Source { return eventSystemSource{} }

func (eventSystemSource) Open() error  { return nil }
func (eventSystemSource) Close() error { return nil }

func (eventSystemSource) Readings() ([]model.Reading, error) {
	var out []model.Reading

	appendReadings := func(kind model.Kind, unit string, cstr *C.char) {
		if cstr == nil {
			return
		}
		defer C.hwmon02_free(cstr)

		s := C.GoString(cstr)
		sc := bufio.NewScanner(strings.NewReader(s))
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			name, valStr, ok := strings.Cut(sc.Text(), ":")
			if !ok {
				continue
			}
			name = strings.TrimSpace(name)
			f, err := strconv.ParseFloat(strings.TrimSpace(valStr), 64)
			if err != nil {
				continue
			}
			f = math.Abs(f)
			f = model.NormalizeValue(kind, name, f)
			out = append(out, model.Reading{
				Kind:    kind,
				Name:    name,
				Unit:    unit,
				Source:  "hid",
				Value:   f,
				KeyOrID: fmt.Sprintf("hid_es:%s", name),
			})
		}
	}

	appendReadings(model.KindAmp, "A", C.hwmon02_get_currents())
	appendReadings(model.KindVolt, "V", C.hwmon02_get_voltages())
	appendReadings(model.KindWatt, "W", C.hwmon02_get_powers())
	appendReadings(model.KindTemp, "°C", C.hwmon02_get_thermals())

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].KeyOrID < out[j].KeyOrID
	})
	return out, nil
}
