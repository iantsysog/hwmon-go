//go:build darwin && cgo && eventsystem

#ifndef __HWMON_HID_02_H__
#define __HWMON_HID_02_H__

char* hwmon02_get_currents();
char* hwmon02_get_voltages();
char* hwmon02_get_powers();
char* hwmon02_get_thermals();
void hwmon02_free(char* s);

#endif
