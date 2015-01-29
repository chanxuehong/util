package random

import (
	"crypto/md5"
	"encoding/hex"
	"sync/atomic"
	"time"
)

func commonRandom(localSalt []byte) []byte {
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

	hashSumArray := md5.Sum(src)
	return hashSumArray[:]
}

// 返回的结果没有经过 hex 编码, 不是可显示的字符串
func NewRawRandom() []byte {
	return commonRandom(localRandomSalt)
}

// 返回的结果已经经过 hex 编码
func NewToken() (token []byte) {
	tokenx := commonRandom(localTokenSalt)
	token = make([]byte, hex.EncodedLen(len(tokenx)))
	hex.Encode(token, tokenx)
	return
}
