package atomic

import (
	"sync/atomic"
	"unsafe"
)

type Bool uint32 // zero value represents false

func LoadBool(addr *Bool) (val bool) {
	n := atomic.LoadUint32((*uint32)(unsafe.Pointer(addr)))
	return n != 0
}

func StoreBool(addr *Bool, val bool) {
	if val {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(addr)), 1)
	} else {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(addr)), 0)
	}
}
