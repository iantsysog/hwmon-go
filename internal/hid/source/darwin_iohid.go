//go:build darwin && cgo && iohid

package source

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/hid/IOHIDLib.h>
#include <stdlib.h>
#include "iohid.h"
*/
import "C"

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"unsafe"

	"github.com/iantsysog/hwmon-go/internal/hid/model"
)

type ioHIDSource struct {
	mu  sync.Mutex
	mgr C.IOHIDManagerRef
}

func defaultIOHIDSource() Source { return &ioHIDSource{} }

func (s *ioHIDSource) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.openLocked()
}

func (s *ioHIDSource) openLocked() error {
	if s.mgr != 0 {
		return nil
	}
	mgr := C.IOHIDManagerCreate(C.kCFAllocatorDefault, C.kIOHIDOptionsTypeNone)
	if mgr == 0 {
		return fmt.Errorf("IOHIDManagerCreate returned NULL")
	}
	C.IOHIDManagerSetDeviceMatching(mgr, 0)
	if ret := C.IOHIDManagerOpen(mgr, C.kIOHIDOptionsTypeNone); ret != C.kIOReturnSuccess {
		C.CFRelease(C.CFTypeRef(mgr))
		return fmt.Errorf("IOHIDManagerOpen failed (ret=0x%x)", int(ret))
	}
	s.mgr = mgr
	return nil
}

func (s *ioHIDSource) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mgr == 0 {
		return nil
	}
	C.IOHIDManagerClose(s.mgr, C.kIOHIDOptionsTypeNone)
	C.CFRelease(C.CFTypeRef(s.mgr))
	s.mgr = 0
	return nil
}

func (s *ioHIDSource) Readings() ([]model.Reading, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.openLocked(); err != nil {
		return nil, err
	}

	var devices *C.IOHIDDeviceRef
	var devCount C.CFIndex
	if ok := C.hwmon01_hid_copy_devices(s.mgr, &devices, &devCount); ok == 0 {
		return nil, fmt.Errorf("failed to enumerate HID devices")
	}
	defer C.hwmon01_hid_free_devices(devices)

	cVendor := C.CString("VendorID")
	cProductID := C.CString("ProductID")
	cProduct := C.CString("Product")
	defer C.free(unsafe.Pointer(cVendor))
	defer C.free(unsafe.Pointer(cProductID))
	defer C.free(unsafe.Pointer(cProduct))

	var out []model.Reading
	for _, dev := range unsafe.Slice(devices, int(devCount)) {
		out = append(out, s.readDevice(dev, cVendor, cProductID, cProduct)...)
	}

	sort.Slice(out, func(i, j int) bool {
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

func (s *ioHIDSource) readDevice(dev C.IOHIDDeviceRef, cVendor, cProductID, cProduct *C.char) []model.Reading {
	var vid, pid int
	_ = C.hwmon01_hid_get_int_property_name(dev, cVendor, (*C.int)(unsafe.Pointer(&vid)))
	_ = C.hwmon01_hid_get_int_property_name(dev, cProductID, (*C.int)(unsafe.Pointer(&pid)))

	var product string
	if c := C.hwmon01_hid_get_string_property_name(dev, cProduct); c != nil {
		product = C.GoString(c)
		C.free(unsafe.Pointer(c))
	}
	if product == "" {
		product = fmt.Sprintf("hid:%04x:%04x", vid, pid)
	}

	var elems *C.IOHIDElementRef
	var elemCount C.CFIndex
	if ok := C.hwmon01_hid_copy_elements(dev, &elems, &elemCount); ok == 0 {
		return nil
	}
	defer C.hwmon01_hid_free_elements(elems)

	out := make([]model.Reading, 0, 16)
	for _, elem := range unsafe.Slice(elems, int(elemCount)) {

		usagePage := uint32(C.IOHIDElementGetUsagePage(elem))
		if usagePage != 0x20 {
			continue
		}

		t := C.IOHIDElementGetType(elem)
		if t != C.kIOHIDElementTypeInput_Misc && t != C.kIOHIDElementTypeInput_Button && t != C.kIOHIDElementTypeInput_Axis {
			continue
		}

		var v C.double
		if C.hwmon01_hid_read_element_double(dev, elem, &v) == 0 {
			continue
		}

		usage := uint32(C.IOHIDElementGetUsage(elem))
		cookie := uint32(C.IOHIDElementGetCookie(elem))

		name := s.elementName(elem)
		if name == "" {
			name = fmt.Sprintf("%s usagePage=0x%x usage=0x%x", product, usagePage, usage)
		}

		kind, unit, ok := model.Classify(model.ElementMeta{Name: name, UsagePage: usagePage, Usage: usage})
		if !ok {
			continue
		}

		value := normalizeIOHIDValue(float64(v), elem)

		out = append(out, model.Reading{
			Kind:    kind,
			Name:    name,
			Unit:    unit,
			Source:  "hid",
			Value:   value,
			KeyOrID: fmt.Sprintf("hid:%04x:%04x:%x:%x:%x", vid, pid, usagePage, usage, cookie),
		})
	}
	return out
}

func (s *ioHIDSource) elementName(elem C.IOHIDElementRef) string {
	cfs := C.IOHIDElementGetName(elem)
	if cfs == 0 {
		return ""
	}
	c := C.hwmon01_cfstring_copy_cstr(cfs)
	if c == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

func normalizeIOHIDValue(raw float64, elem C.IOHIDElementRef) float64 {
	value := raw
	var u, exp C.int
	if C.hwmon01_hid_element_unit(elem, &u, &exp) != 0 {
		if exp != 0 {
			value = value * math.Pow10(int(exp))
		}
		_ = u
	}
	return value
}
