package atomic

import (
	"sync/atomic"
	"unsafe"
)

type Bool uint32 // zero value represents false

func (b *Bool) Load() (val bool) {
	n := atomic.LoadUint32((*uint32)(unsafe.Pointer(b)))
	return n != 0
}

func (b *Bool) Store(val bool) {
	if val {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(b)), 1)
	} else {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(b)), 0)
	}
}
