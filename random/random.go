package random

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"encoding/hex"
	mathrand "math/rand"
	"sync"
	"time"
)

const (
	randomSaltLen            = 45   // see NewRandom(), 8+2+45 == 55, md5 签名性能最好的前提下最大的数据块
	randomSaltUpdateInterval = 3600 // seconds
)

var (
	randomSalt               []byte = make([]byte, randomSaltLen)
	randomMutex              sync.Mutex
	randomSaltLastUpdateTime int64  = time.Now().Unix()
	randomClockSequence      uint32 = NewRandomUint32()
)

func init() {
	ReadRandomBytes(randomSalt)
}

// NewRandom 返回一个随机字节数组.
//  NOTE: 返回的是原始数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符
func NewRandom() (ret [16]byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()

	randomMutex.Lock() // Lock
	if timeNowUnix >= randomSaltLastUpdateTime+randomSaltUpdateInterval {
		randomSaltLastUpdateTime = timeNowUnix
		ReadRandomBytes(randomSalt)

		randomMutex.Unlock() // Unlock
		copy(ret[:], randomSalt)
		return
	}
	randomClockSequence++
	clockSequence := randomClockSequence
	randomMutex.Unlock() // Unlock

	var src [8 + 2 + randomSaltLen]byte // 8+2+45 == 55
	timeNowUnixNano := timeNow.UnixNano()
	src[0] = byte(timeNowUnixNano >> 56)
	src[1] = byte(timeNowUnixNano >> 48)
	src[2] = byte(timeNowUnixNano >> 40)
	src[3] = byte(timeNowUnixNano >> 32)
	src[4] = byte(timeNowUnixNano >> 24)
	src[5] = byte(timeNowUnixNano >> 16)
	src[6] = byte(timeNowUnixNano >> 8)
	src[7] = byte(timeNowUnixNano)
	src[8] = byte(clockSequence >> 8)
	src[9] = byte(clockSequence)
	copy(src[10:], randomSalt)

	ret = md5.Sum(src[:])
	return
}

// NewRandomEx 返回一个32字节的随机数.
//  NOTE: 返回的结果经过了 hex 编码.
func NewRandomEx() (ret []byte) {
	rd := NewRandom()
	ret = make([]byte, hex.EncodedLen(len(rd)))
	hex.Encode(ret, rd[:])
	return
}

var globalMathRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

// 读取随机字节到 p []byte 里面.
func ReadRandomBytes(p []byte) {
	if len(p) <= 0 {
		return
	}

	// get from crypto/rand
	if _, err := cryptorand.Read(p); err == nil {
		return
	}

	// get from math/rand
	timeNowUnixNano := time.Now().UnixNano()
	for len(p) > 0 {
		switch n := globalMathRand.Int63() ^ timeNowUnixNano; len(p) {
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
