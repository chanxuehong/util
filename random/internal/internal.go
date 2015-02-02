// random 内部包, 请不要使用, 因为在以后的go版本这个 package 也许不能调用了.
package internal

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	mathRand "math/rand"
	"time"
)

var globalMathRand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

// 读取随机的字节到 p 指向的 []byte 里面.
func ReadRandomBytes(p []byte) {
	if len(p) <= 0 {
		return
	}

	// get from crypto/rand
	if _, err := cryptoRand.Read(p); err == nil {
		return
	}

	// get from math/rand
	for len(p) > 0 {
		n := globalMathRand.Int63()

		switch len(p) {
		case 8:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			p[6] = byte(n >> 8)
			p[7] = byte(n)
			return
		case 4:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			return
		case 1:
			p[0] = byte(n >> 56)
			return
		case 2:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			return
		case 3:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			return
		case 5:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			return
		case 6:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			return
		case 7:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			p[6] = byte(n >> 8)
			return
		default: // len(p) > 8
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			p[6] = byte(n >> 8)
			p[7] = byte(n)
			p = p[8:]
		}
	}
}

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
