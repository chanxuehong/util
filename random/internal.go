package random

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	mathRand "math/rand"
	"time"
)

// 全局的 math/rand.Rand
//  NOTE:
//  使用 var 初始化, 其他直接或简洁依赖 readRandomBytes 的全局变量都是通过 init() 来初始化.
var globalMathRand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

// 读取随机的字节到 p []byte 里面.
func readRandomBytes(p []byte) {
	if len(p) <= 0 {
		return
	}

	// get from crypto/rand
	if _, err := cryptoRand.Read(p); err == nil {
		return
	}

	// get from math/rand
	for buf := p; len(buf) > 0; {
		n := globalMathRand.Int63()

		switch len(buf) {
		case 8:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			buf[4] = byte(n >> 24)
			buf[5] = byte(n >> 16)
			buf[6] = byte(n >> 8)
			buf[7] = byte(n)
			return
		case 4:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			return
		case 1:
			buf[0] = byte(n >> 56)
			return
		case 2:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			return
		case 3:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			return
		case 5:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			buf[4] = byte(n >> 24)
			return
		case 6:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			buf[4] = byte(n >> 24)
			buf[5] = byte(n >> 16)
			return
		case 7:
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			buf[4] = byte(n >> 24)
			buf[5] = byte(n >> 16)
			buf[6] = byte(n >> 8)
			return
		default: // len(buf) > 8
			buf[0] = byte(n >> 56)
			buf[1] = byte(n >> 48)
			buf[2] = byte(n >> 40)
			buf[3] = byte(n >> 32)
			buf[4] = byte(n >> 24)
			buf[5] = byte(n >> 16)
			buf[6] = byte(n >> 8)
			buf[7] = byte(n)
			buf = buf[8:]
		}
	}
}

// 获取一个随机的 uint32 整数.
func newRandomUint32() uint32 {
	var x [4]byte
	readRandomBytes(x[:])
	return binary.BigEndian.Uint32(x[:])
}

// 获取一个随机的 uint64 整数.
func newRandomUint64() uint64 {
	var x [8]byte
	readRandomBytes(x[:])
	return binary.BigEndian.Uint64(x[:])
}
