package random

import (
	"encoding/binary"
)

// 获取一个随机的 uint32 整数.
func NewRandomUint32() uint32 {
	var x [4]byte
	ReadRandomBytes(x[:])
	return binary.BigEndian.Uint32(x[:])
}

// 获取一个随机的 uint64 整数.
func NewRandomUint64() uint64 {
	var x [8]byte
	ReadRandomBytes(x[:])
	return binary.BigEndian.Uint64(x[:])
}
