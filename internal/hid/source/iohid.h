//go:build darwin && cgo && iohid

#ifndef __HWMON_HID_01_H__
#define __HWMON_HID_01_H__

#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/hid/IOHIDLib.h>

char* hwmon01_cfstring_copy_cstr(CFStringRef s);

int hwmon01_hid_copy_devices(IOHIDManagerRef mgr, IOHIDDeviceRef** outDevices, CFIndex* outCount);
void hwmon01_hid_free_devices(IOHIDDeviceRef* devices);

int hwmon01_hid_copy_elements(IOHIDDeviceRef dev, IOHIDElementRef** outElems, CFIndex* outCount);
void hwmon01_hid_free_elements(IOHIDElementRef* elems);

int hwmon01_hid_get_int_property_name(IOHIDDeviceRef dev, const char* keyName, int* out);
char* hwmon01_hid_get_string_property_name(IOHIDDeviceRef dev, const char* keyName);

int hwmon01_hid_read_element_double(IOHIDDeviceRef dev, IOHIDElementRef elem, double* out);
int hwmon01_hid_element_unit(IOHIDElementRef elem, int* outUnit, int* outExp);

#endif
