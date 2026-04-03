//go:build darwin && cgo

package smc

/*
#cgo LDFLAGS: -framework IOKit
#include <stdlib.h>
#include "smc.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type AppleConnection struct {
	conn    C.io_connect_t
	opened  atomic.Bool
	connMu  sync.RWMutex
	cacheMu sync.RWMutex
	cache   map[string]C.SMCKeyData_keyInfo_t
}

func New() Connection {
	return &AppleConnection{
		cache: make(map[string]C.SMCKeyData_keyInfo_t),
	}
}

func (c *AppleConnection) Open() error {
	c.connMu.Lock()
	defer c.connMu.Unlock()

	if c.opened.Load() {
		return nil
	}
	if ret := int(C.SMCOpen(&c.conn)); ret != 0 {
		c.conn = 0
		return &SMCError{Op: "open connection", Code: ret}
	}
	c.opened.Store(true)
	return nil
}

func (c *AppleConnection) Close() error {
	c.connMu.Lock()
	defer c.connMu.Unlock()

	if !c.opened.Load() {
		return nil
	}

	if ret := int(C.SMCClose(c.conn)); ret != 0 {
		return &SMCError{Op: "close connection", Code: ret}
	}
	c.opened.Store(false)
	c.cacheMu.Lock()
	clear(c.cache)
	c.cacheMu.Unlock()
	c.conn = 0
	return nil
}

func (c *AppleConnection) Write(key string, val []byte) error {
	if len(key) != 4 {
		return ErrInvalidKey
	}
	if !c.opened.Load() {
		return ErrConnectionClosed
	}
	if len(val) > 32 {
		return fmt.Errorf("%w: %d", ErrInvalidDataSize, len(val))
	}

	c.connMu.RLock()
	defer c.connMu.RUnlock()
	if !c.opened.Load() {
		return ErrConnectionClosed
	}

	ckey := C.CString(key)
	cval := C.CBytes(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(cval)

	if ret := int(C.SMCWriteSimple(ckey, (*C.uchar)(cval), C.int(len(val)), c.conn)); ret != 0 {
		return &SMCError{Op: "write key", Key: key, Code: ret}
	}

	return nil
}

func (c *AppleConnection) Read(key string) (SMCVal, error) {
	if len(key) != 4 {
		return SMCVal{}, ErrInvalidKey
	}
	if !c.opened.Load() {
		return SMCVal{}, ErrConnectionClosed
	}

	c.connMu.RLock()
	defer c.connMu.RUnlock()
	if !c.opened.Load() {
		return SMCVal{}, ErrConnectionClosed
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	keyInfo, err := c.getKeyInfo(key)
	if err != nil {
		return SMCVal{}, err
	}

	v := C.SMCVal_t{}

	if ret := int(C.SMCReadKeyWithInfo2(ckey, &keyInfo, &v, c.conn)); ret != 0 {
		return SMCVal{}, &SMCError{Op: "read key", Key: key, Code: ret}
	}

	if v.dataSize == 0 {
		return SMCVal{}, ErrNoDataForKey
	}
	if v.dataSize > 32 {
		return SMCVal{}, fmt.Errorf("%w: %d", ErrInvalidDataSize, v.dataSize)
	}

	bytes := C.GoBytes(unsafe.Pointer(&v.bytes[0]), C.int(v.dataSize))

	val := SMCVal{
		Key:      key,
		DataType: fourCCFromUint32(uint32(keyInfo.dataType)),
		Bytes:    bytes,
	}

	return val, nil
}

func (c *AppleConnection) getKeyInfo(key string) (C.SMCKeyData_keyInfo_t, error) {
	c.cacheMu.RLock()
	keyInfo, ok := c.cache[key]
	c.cacheMu.RUnlock()
	if ok {
		return keyInfo, nil
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	if ret := int(C.SMCReadKeyInfo2(ckey, &keyInfo, c.conn)); ret != 0 {
		return C.SMCKeyData_keyInfo_t{}, &SMCError{Op: "read key info", Key: key, Code: ret}
	}

	c.cacheMu.Lock()
	if c.cache == nil {
		c.cache = make(map[string]C.SMCKeyData_keyInfo_t)
	}
	c.cache[key] = keyInfo
	c.cacheMu.Unlock()

	return keyInfo, nil
}

func fourCCFromUint32(v uint32) string {
	return string([]byte{
		byte((v >> 24) & 0xff),
		byte((v >> 16) & 0xff),
		byte((v >> 8) & 0xff),
		byte(v & 0xff),
	})
}
