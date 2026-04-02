#ifndef __SMC_H__
#define __SMC_H__

#include <stdint.h>
#include <IOKit/IOKitLib.h>

#define OP_NONE 0
#define OP_LIST 1
#define OP_READ 2
#define OP_READ_FAN 3
#define OP_WRITE 4
#define OP_READ_TEMPS 5

#define KERNEL_INDEX_SMC 2

#define SMC_CMD_READ_BYTES 5
#define SMC_CMD_WRITE_BYTES 6
#define SMC_CMD_READ_INDEX 8
#define SMC_CMD_READ_KEYINFO 9
#define SMC_CMD_READ_PLIMIT 11
#define SMC_CMD_READ_VERS 12

typedef struct {
    uint8_t major;
    uint8_t minor;
    uint8_t build;
    uint8_t reserved[1];
    uint16_t release;
} SMCKeyData_vers_t;

typedef struct {
    uint16_t version;
    uint16_t length;
    uint32_t cpuPLimit;
    uint32_t gpuPLimit;
    uint32_t memPLimit;
} SMCKeyData_pLimitData_t;

typedef struct {
    uint32_t dataSize;
    uint32_t dataType;
    uint8_t dataAttributes;
} SMCKeyData_keyInfo_t;

typedef unsigned char SMCBytes_t[32];

typedef struct {
    uint32_t key;
    SMCKeyData_vers_t vers;
    SMCKeyData_pLimitData_t pLimitData;
    SMCKeyData_keyInfo_t keyInfo;
    uint8_t result;
    uint8_t status;
    uint8_t data8;
    uint32_t data32;
    SMCBytes_t bytes;
} SMCKeyData_t;

typedef char UInt32Char_t[5];

typedef struct {
    UInt32Char_t key;
    uint32_t dataSize;
    UInt32Char_t dataType;
    SMCBytes_t bytes;
} SMCVal_t;

kern_return_t SMCReadKey(UInt32Char_t key, SMCVal_t* val);
kern_return_t SMCWriteSimple(const char* key, const unsigned char* bytes, int len, io_connect_t conn);

kern_return_t SMCOpen(io_connect_t* conn);
kern_return_t SMCClose(io_connect_t conn);
kern_return_t SMCReadKey2(const UInt32Char_t key, SMCVal_t* val, io_connect_t conn);
kern_return_t SMCReadKeyInfo2(const UInt32Char_t key, SMCKeyData_keyInfo_t* keyInfo, io_connect_t conn);
kern_return_t SMCReadKeyWithInfo2(const UInt32Char_t key, const SMCKeyData_keyInfo_t* keyInfo, SMCVal_t* val, io_connect_t conn);

#endif
