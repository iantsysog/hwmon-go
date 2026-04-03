//go:build darwin && cgo && iohid

#include "iohid.h"
#include <stdlib.h>

char* hwmon01_cfstring_copy_cstr(CFStringRef s) {
    if (s == NULL) return NULL;
    CFIndex len = CFStringGetLength(s);
    CFIndex maxSize = CFStringGetMaximumSizeForEncoding(len, kCFStringEncodingUTF8) + 1;
    char *buf = (char*)malloc((size_t)maxSize);
    if (buf == NULL) return NULL;
    if (!CFStringGetCString(s, buf, maxSize, kCFStringEncodingUTF8)) {
        free(buf);
        return NULL;
    }
    return buf;
}

int hwmon01_hid_copy_devices(IOHIDManagerRef mgr, IOHIDDeviceRef** outDevices, CFIndex* outCount) {
    if (outDevices == NULL || outCount == NULL) return 0;
    *outDevices = NULL;
    *outCount = 0;
    if (mgr == NULL) return 0;

    CFSetRef set = IOHIDManagerCopyDevices(mgr);
    if (set == NULL) return 1;
    CFIndex count = CFSetGetCount(set);
    *outCount = count;
    if (count <= 0) {
        CFRelease(set);
        return 1;
    }

    IOHIDDeviceRef* arr = (IOHIDDeviceRef*)calloc((size_t)count, sizeof(IOHIDDeviceRef));
    if (arr == NULL) {
        CFRelease(set);
        return 0;
    }
    CFSetGetValues(set, (const void**)arr);
    *outDevices = arr;
    CFRelease(set);
    return 1;
}

void hwmon01_hid_free_devices(IOHIDDeviceRef* devices) {
    free(devices);
}

int hwmon01_hid_copy_elements(IOHIDDeviceRef dev, IOHIDElementRef** outElems, CFIndex* outCount) {
    if (outElems == NULL || outCount == NULL) return 0;
    *outElems = NULL;
    *outCount = 0;
    if (dev == NULL) return 0;

    CFArrayRef arr = IOHIDDeviceCopyMatchingElements(dev, NULL, kIOHIDOptionsTypeNone);
    if (arr == NULL) return 1;
    CFIndex count = CFArrayGetCount(arr);
    *outCount = count;
    if (count <= 0) {
        CFRelease(arr);
        return 1;
    }

    IOHIDElementRef* elems = (IOHIDElementRef*)calloc((size_t)count, sizeof(IOHIDElementRef));
    if (elems == NULL) {
        CFRelease(arr);
        return 0;
    }
    for (CFIndex i = 0; i < count; i++) {
        elems[i] = (IOHIDElementRef)CFArrayGetValueAtIndex(arr, i);
    }
    *outElems = elems;
    CFRelease(arr);
    return 1;
}

void hwmon01_hid_free_elements(IOHIDElementRef* elems) {
    free(elems);
}

int hwmon01_hid_get_int_property_name(IOHIDDeviceRef dev, const char* keyName, int* out) {
    if (out == NULL) return 0;
    *out = 0;
    if (dev == NULL || keyName == NULL) return 0;
    CFStringRef key = CFStringCreateWithCString(kCFAllocatorDefault, keyName, kCFStringEncodingUTF8);
    if (key == NULL) return 0;
    CFTypeRef ref = IOHIDDeviceGetProperty(dev, key);
    CFRelease(key);
    if (ref == NULL) return 0;
    if (CFGetTypeID(ref) != CFNumberGetTypeID()) return 0;
    int v = 0;
    if (!CFNumberGetValue((CFNumberRef)ref, kCFNumberIntType, &v)) return 0;
    *out = v;
    return 1;
}

char* hwmon01_hid_get_string_property_name(IOHIDDeviceRef dev, const char* keyName) {
    if (dev == NULL || keyName == NULL) return NULL;
    CFStringRef key = CFStringCreateWithCString(kCFAllocatorDefault, keyName, kCFStringEncodingUTF8);
    if (key == NULL) return NULL;
    CFTypeRef ref = IOHIDDeviceGetProperty(dev, key);
    CFRelease(key);
    if (ref == NULL) return NULL;
    if (CFGetTypeID(ref) != CFStringGetTypeID()) return NULL;
    return hwmon01_cfstring_copy_cstr((CFStringRef)ref);
}

int hwmon01_hid_read_element_double(IOHIDDeviceRef dev, IOHIDElementRef elem, double* out) {
    if (out == NULL) return 0;
    *out = 0;
    if (dev == NULL || elem == NULL) return 0;
    IOHIDValueRef value = NULL;
    IOReturn ret = IOHIDDeviceGetValue(dev, elem, &value);
    if (ret != kIOReturnSuccess || value == NULL) return 0;
    CFIndex i = IOHIDValueGetIntegerValue(value);
    *out = (double)i;
    return 1;
}

int hwmon01_hid_element_unit(IOHIDElementRef elem, int* outUnit, int* outExp) {
    if (outUnit == NULL || outExp == NULL) return 0;
    *outUnit = 0;
    *outExp = 0;
    if (elem == NULL) return 0;
    *outUnit = (int)IOHIDElementGetUnit(elem);
    *outExp = (int)IOHIDElementGetUnitExponent(elem);
    return 1;
}
