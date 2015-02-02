package random

import (
	"crypto/md5"
	"encoding/hex"
	"sync/atomic"
	"time"
)

const randomSaltLen = 45

var (
	randomSalt          = make([]byte, randomSaltLen)
	randomClockSequence uint32
)

func init() {
	readRandomBytes(randomSalt)
	randomClockSequence = newRandomUint32()
}

// NewRandom 返回一个随机字节数组.
//  NOTE: 返回的是原始数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符
func NewRandom() [16]byte {
	var src [8 + 2 + randomSaltLen]byte // 8+2+45 == 55

	nowUnixNano := time.Now().UnixNano()
	src[0] = byte(nowUnixNano >> 56)
	src[1] = byte(nowUnixNano >> 48)
	src[2] = byte(nowUnixNano >> 40)
	src[3] = byte(nowUnixNano >> 32)
	src[4] = byte(nowUnixNano >> 24)
	src[5] = byte(nowUnixNano >> 16)
	src[6] = byte(nowUnixNano >> 8)
	src[7] = byte(nowUnixNano)

	seq := atomic.AddUint32(&randomClockSequence, 1)
	src[8] = byte(seq >> 8)
	src[9] = byte(seq)

	copy(src[10:], randomSalt)

	return md5.Sum(src[:])
}

// NewToken 返回一个32字节的随机数.
//  NOTE: 返回的结果经过了 hex 编码.
func NewToken() (token []byte) {
	random := NewRandom()
	token = make([]byte, hex.EncodedLen(len(random)))
	hex.Encode(token, random[:])
	return
}
