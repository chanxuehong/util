package random

import (
	"crypto/md5"
	"encoding/hex"
	"sync/atomic"
	"time"
)

func commonRandom(localSalt []byte) (random [16]byte) {
	src := make([]byte, 8+2+localSaltLen) // nowNanosecond + seq + localSalt

	nowNanosecond := time.Now().UnixNano()
	src[0] = byte(nowNanosecond >> 56)
	src[1] = byte(nowNanosecond >> 48)
	src[2] = byte(nowNanosecond >> 40)
	src[3] = byte(nowNanosecond >> 32)
	src[4] = byte(nowNanosecond >> 24)
	src[5] = byte(nowNanosecond >> 16)
	src[6] = byte(nowNanosecond >> 8)
	src[7] = byte(nowNanosecond)

	seq := atomic.AddUint32(&randomClockSequence, 1)
	src[8] = byte(seq >> 8)
	src[9] = byte(seq)

	copy(src[10:], localSalt)

	random = md5.Sum(src)
	return
}

// NewRandom 返回一个随机数.
//  NOTE: 返回的是原始数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符
func NewRandom() [16]byte {
	return commonRandom(localRandomSalt)
}

// NewToken 返回一个32字节的随机数.
//  NOTE: 返回的结果已经经过 hex 编码.
func NewToken() (token []byte) {
	random := commonRandom(localTokenSalt)
	token = make([]byte, hex.EncodedLen(len(random)))
	hex.Encode(token, random[:])
	return
}
