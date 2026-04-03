//go:build darwin && cgo && eventsystem

#import <Foundation/Foundation.h>
#import <IOKit/hidsystem/IOHIDEventSystemClient.h>
#include <stdlib.h>
#include <string.h>

typedef struct __IOHIDEvent *IOHIDEventRef;
typedef struct __IOHIDServiceClient *IOHIDServiceClientRef;
typedef double IOHIDFloat;

IOHIDEventSystemClientRef IOHIDEventSystemClientCreate(CFAllocatorRef allocator);
int IOHIDEventSystemClientSetMatching(IOHIDEventSystemClientRef client, CFDictionaryRef match);
CFArrayRef IOHIDEventSystemClientCopyServices(IOHIDEventSystemClientRef client);

IOHIDEventRef IOHIDServiceClientCopyEvent(IOHIDServiceClientRef, int64_t, int32_t, int64_t);
CFStringRef IOHIDServiceClientCopyProperty(IOHIDServiceClientRef service, CFStringRef property);
IOHIDFloat IOHIDEventGetFloatValue(IOHIDEventRef event, int32_t field);

static NSDictionary* hwmon02_matching(int page, int usage) {
    return @{
        @"PrimaryUsagePage" : [NSNumber numberWithInt:page],
        @"PrimaryUsage" : [NSNumber numberWithInt:usage],
    };
}

static NSArray* hwmon02_getProductNames(NSDictionary* sensors) {
    IOHIDEventSystemClientRef system = IOHIDEventSystemClientCreate(kCFAllocatorDefault);
    if (!system) return [[NSArray alloc] init];

    IOHIDEventSystemClientSetMatching(system, (__bridge CFDictionaryRef)sensors);
    CFArrayRef matchingsrvsRef = IOHIDEventSystemClientCopyServices(system);
    NSArray* matchingsrvs = matchingsrvsRef ? (__bridge NSArray*)matchingsrvsRef : @[];

    long count = [matchingsrvs count];
    NSMutableArray* array = [[NSMutableArray alloc] init];
    for (int i = 0; i < count; i++) {
        IOHIDServiceClientRef sc = (IOHIDServiceClientRef)matchingsrvs[i];
        NSString* name = (NSString*)IOHIDServiceClientCopyProperty(sc, (__bridge CFStringRef)@"Product");
        if (name) {
            [array addObject:name];
            [name release];
        } else {
            [array addObject:@"noname"];
        }
    }

    if (matchingsrvsRef) CFRelease(matchingsrvsRef);
    CFRelease(system);
    return array;
}

#define IOHIDEventFieldBase(type) (type << 16)
#define kIOHIDEventTypeTemperature 15
#define kIOHIDEventTypePower 25

static NSArray* hwmon02_getPowerValues(NSDictionary* sensors) {
    IOHIDEventSystemClientRef system = IOHIDEventSystemClientCreate(kCFAllocatorDefault);
    if (!system) return [[NSArray alloc] init];

    IOHIDEventSystemClientSetMatching(system, (__bridge CFDictionaryRef)sensors);
    CFArrayRef matchingsrvsRef = IOHIDEventSystemClientCopyServices(system);
    NSArray* matchingsrvs = matchingsrvsRef ? (__bridge NSArray*)matchingsrvsRef : @[];

    long count = [matchingsrvs count];
    NSMutableArray* array = [[NSMutableArray alloc] init];
    for (int i = 0; i < count; i++) {
        IOHIDServiceClientRef sc = (IOHIDServiceClientRef)matchingsrvs[i];
        IOHIDEventRef event = IOHIDServiceClientCopyEvent(sc, kIOHIDEventTypePower, 0, 0);

        double val = 0.0;
        if (event != 0) {
            val = IOHIDEventGetFloatValue(event, IOHIDEventFieldBase(kIOHIDEventTypePower)) / 1000.0;
            CFRelease(event);
        }

        [array addObject:[NSNumber numberWithDouble:fabs(val)]];
    }

    if (matchingsrvsRef) CFRelease(matchingsrvsRef);
    CFRelease(system);
    return array;
}

static NSArray* hwmon02_getThermalValues(NSDictionary* sensors) {
    IOHIDEventSystemClientRef system = IOHIDEventSystemClientCreate(kCFAllocatorDefault);
    if (!system) return [[NSArray alloc] init];

    IOHIDEventSystemClientSetMatching(system, (__bridge CFDictionaryRef)sensors);
    CFArrayRef matchingsrvsRef = IOHIDEventSystemClientCopyServices(system);
    NSArray* matchingsrvs = matchingsrvsRef ? (__bridge NSArray*)matchingsrvsRef : @[];

    long count = [matchingsrvs count];
    NSMutableArray* array = [[NSMutableArray alloc] init];
    for (int i = 0; i < count; i++) {
        IOHIDServiceClientRef sc = (IOHIDServiceClientRef)matchingsrvs[i];
        IOHIDEventRef event = IOHIDServiceClientCopyEvent(sc, kIOHIDEventTypeTemperature, 0, 0);

        double val = 0.0;
        if (event != 0) {
            val = IOHIDEventGetFloatValue(event, IOHIDEventFieldBase(kIOHIDEventTypeTemperature));
            CFRelease(event);
        }

        [array addObject:[NSNumber numberWithDouble:fabs(val)]];
    }

    if (matchingsrvsRef) CFRelease(matchingsrvsRef);
    CFRelease(system);
    return array;
}

static char* hwmon02_dump(NSArray* names, NSArray* values) {
    NSMutableString* valueString = [[NSMutableString alloc] init];
    int count = (int)MIN([names count], [values count]);
    for (int i = 0; i < count; i++) {
        @autoreleasepool {
            NSString* name = names[i];
            double value = [values[i] doubleValue];
            [valueString appendFormat:@"%s:%lf\n", [name UTF8String], value];
        }
    }

    const char* utf8 = valueString ? [valueString UTF8String] : "";
    char* finalStr = strdup(utf8 ? utf8 : "");
    [valueString release];
    return finalStr;
}

char* hwmon02_get_currents() {
    @autoreleasepool {
        NSDictionary* sensors = hwmon02_matching(0xff08, 2);
        NSArray* names = hwmon02_getProductNames(sensors);
        NSArray* vals = hwmon02_getPowerValues(sensors);
        char* out = hwmon02_dump(names, vals);
        CFRelease(names);
        CFRelease(vals);
        return out;
    }
}

char* hwmon02_get_voltages() {
    @autoreleasepool {
        NSDictionary* sensors = hwmon02_matching(0xff08, 3);
        NSArray* names = hwmon02_getProductNames(sensors);
        NSArray* vals = hwmon02_getPowerValues(sensors);
        char* out = hwmon02_dump(names, vals);
        CFRelease(names);
        CFRelease(vals);
        return out;
    }
}

char* hwmon02_get_powers() {
    @autoreleasepool {
        NSDictionary* sensors = hwmon02_matching(0xff08, 1);
        NSArray* names = hwmon02_getProductNames(sensors);
        NSArray* vals = hwmon02_getPowerValues(sensors);
        char* out = hwmon02_dump(names, vals);
        CFRelease(names);
        CFRelease(vals);
        return out;
    }
}

char* hwmon02_get_thermals() {
    @autoreleasepool {
        NSDictionary* sensors = hwmon02_matching(0xff00, 5);
        NSArray* names = hwmon02_getProductNames(sensors);
        NSArray* vals = hwmon02_getThermalValues(sensors);
        char* out = hwmon02_dump(names, vals);
        CFRelease(names);
        CFRelease(vals);
        return out;
    }
}

void hwmon02_free(char* s) {
    free(s);
}
