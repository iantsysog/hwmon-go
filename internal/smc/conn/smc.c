//go:build darwin && cgo

#include "smc.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

kern_return_t SMCCall2(int index, SMCKeyData_t* inputStructure, SMCKeyData_t* outputStructure, io_connect_t conn);

#pragma mark C Helpers

static uint32_t SMCFourCCToUInt32(const char str[static 4])
{
    return ((uint32_t)(uint8_t)str[0] << 24) |
        ((uint32_t)(uint8_t)str[1] << 16) |
        ((uint32_t)(uint8_t)str[2] << 8) |
        ((uint32_t)(uint8_t)str[3]);
}

static void SMCUInt32ToFourCC(char out[static 5], uint32_t val)
{
    out[0] = (char)((val >> 24) & 0xff);
    out[1] = (char)((val >> 16) & 0xff);
    out[2] = (char)((val >> 8) & 0xff);
    out[3] = (char)(val & 0xff);
    out[4] = '\0';
}

static bool SMCValidateKey(const char* key)
{
    return key != NULL && strlen(key) == 4;
}

static void SMCCopyKey4(UInt32Char_t dst, const char src[static 4])
{
    memcpy(dst, src, 4);
    dst[4] = '\0';
}

#pragma mark Shared SMC functions

kern_return_t SMCOpen(io_connect_t* conn)
{
    if (conn == NULL) {
        return kIOReturnBadArgument;
    }
    *conn = 0;

    mach_port_t masterPort = MACH_PORT_NULL;
    kern_return_t result = IOMainPort(MACH_PORT_NULL, &masterPort);
    if (result != kIOReturnSuccess) {
        return result;
    }

    io_iterator_t iterator = IO_OBJECT_NULL;
    CFMutableDictionaryRef matchingDictionary = IOServiceMatching("AppleSMC");
    if (matchingDictionary == NULL) {
        return kIOReturnNoMemory;
    }
    result = IOServiceGetMatchingServices(masterPort, matchingDictionary, &iterator);
    if (result != kIOReturnSuccess) {
        return result;
    }

    io_object_t device = IOIteratorNext(iterator);
    IOObjectRelease(iterator);
    if (device == 0) {
        return kIOReturnNotFound;
    }

    result = IOServiceOpen(device, mach_task_self(), 0, conn);
    IOObjectRelease(device);
    if (result != kIOReturnSuccess) {
        *conn = 0;
        return result;
    }

    return kIOReturnSuccess;
}

kern_return_t SMCClose(io_connect_t conn)
{
    if (conn == 0) {
        return kIOReturnBadArgument;
    }
    return IOServiceClose(conn);
}

kern_return_t SMCCall2(int index, SMCKeyData_t* inputStructure, SMCKeyData_t* outputStructure, io_connect_t conn)
{
    if (inputStructure == NULL || outputStructure == NULL || conn == 0) {
        return kIOReturnBadArgument;
    }
    size_t structureOutputSize = sizeof(*outputStructure);
    return IOConnectCallStructMethod(conn, index, inputStructure, sizeof(*inputStructure), outputStructure, &structureOutputSize);
}

kern_return_t SMCReadKeyInfo2(const UInt32Char_t key, SMCKeyData_keyInfo_t* keyInfo, io_connect_t conn)
{
    if (keyInfo == NULL || conn == 0 || !SMCValidateKey(key)) {
        return kIOReturnBadArgument;
    }

    SMCKeyData_t inputStructure = { 0 };
    SMCKeyData_t outputStructure = { 0 };
    inputStructure.key = SMCFourCCToUInt32(key);
    inputStructure.data8 = SMC_CMD_READ_KEYINFO;

    kern_return_t result = SMCCall2(KERNEL_INDEX_SMC, &inputStructure, &outputStructure, conn);
    if (result == kIOReturnSuccess) {
        *keyInfo = outputStructure.keyInfo;
    }
    return result;
}

kern_return_t SMCReadKeyWithInfo2(const UInt32Char_t key, const SMCKeyData_keyInfo_t* keyInfo, SMCVal_t* val, io_connect_t conn)
{
    if (val == NULL || keyInfo == NULL || conn == 0 || !SMCValidateKey(key)) {
        return kIOReturnBadArgument;
    }

    SMCKeyData_t inputStructure = { 0 };
    SMCKeyData_t outputStructure = { 0 };
    memset(val, 0, sizeof(*val));

    inputStructure.key = SMCFourCCToUInt32(key);
    SMCCopyKey4(val->key, key);
    val->dataSize = keyInfo->dataSize;
    SMCUInt32ToFourCC(val->dataType, keyInfo->dataType);
    inputStructure.keyInfo.dataSize = val->dataSize;
    inputStructure.data8 = SMC_CMD_READ_BYTES;

    kern_return_t result = SMCCall2(KERNEL_INDEX_SMC, &inputStructure, &outputStructure, conn);
    if (result != kIOReturnSuccess) {
        return result;
    }

    memcpy(val->bytes, outputStructure.bytes, sizeof(val->bytes));

    return kIOReturnSuccess;
}

kern_return_t SMCReadKey2(const UInt32Char_t key, SMCVal_t* val, io_connect_t conn)
{
    SMCKeyData_keyInfo_t keyInfo = { 0 };
    kern_return_t result = SMCReadKeyInfo2(key, &keyInfo, conn);
    if (result != kIOReturnSuccess) {
        return result;
    }

    return SMCReadKeyWithInfo2(key, &keyInfo, val, conn);
}

static kern_return_t SMCWriteKey2(const SMCVal_t* writeVal, io_connect_t conn)
{
    if (writeVal == NULL || conn == 0 || !SMCValidateKey(writeVal->key)) {
        return kIOReturnBadArgument;
    }

    SMCVal_t readVal = { 0 };
    kern_return_t result = SMCReadKey2(writeVal->key, &readVal, conn);
    if (result != kIOReturnSuccess) {
        return result;
    }

    if (readVal.dataSize != writeVal->dataSize) {
        return kIOReturnBadArgument;
    }
    if (writeVal->dataSize > sizeof(writeVal->bytes)) {
        return kIOReturnBadArgument;
    }

    SMCKeyData_t inputStructure = { 0 };
    SMCKeyData_t outputStructure = { 0 };

    inputStructure.key = SMCFourCCToUInt32(writeVal->key);
    inputStructure.data8 = SMC_CMD_WRITE_BYTES;
    inputStructure.keyInfo.dataSize = writeVal->dataSize;
    memcpy(inputStructure.bytes, writeVal->bytes, writeVal->dataSize);
    result = SMCCall2(KERNEL_INDEX_SMC, &inputStructure, &outputStructure, conn);

    return result;
}

kern_return_t SMCWriteSimple(const char* key, const unsigned char* bytes, int len, io_connect_t conn)
{
    if (!SMCValidateKey(key) || bytes == NULL || conn == 0) {
        return kIOReturnBadArgument;
    }

    if (len <= 0 || (size_t)len > sizeof(((SMCVal_t*)0)->bytes)) {
        return kIOReturnBadArgument;
    }

    SMCVal_t val = { 0 };
    SMCCopyKey4(val.key, key);
    val.dataSize = (uint32_t)len;
    memcpy(val.bytes, bytes, (size_t)len);

    return SMCWriteKey2(&val, conn);
}

kern_return_t SMCReadKey(UInt32Char_t key, SMCVal_t* val)
{
    if (val == NULL || !SMCValidateKey(key)) {
        return kIOReturnBadArgument;
    }

    io_connect_t conn = 0;
    kern_return_t result = SMCOpen(&conn);
    if (result != kIOReturnSuccess) {
        return result;
    }

    result = SMCReadKey2(key, val, conn);

    kern_return_t closeResult = SMCClose(conn);
    if (result == kIOReturnSuccess) {
        result = closeResult;
    }
    return result;
}
