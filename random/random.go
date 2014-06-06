package random

import (
	"crypto/md5"
	"encoding/hex"
	"sync/atomic"
	"time"
)

func commonRandom(localSalt []byte, seq uint32) []byte {
	// nowNanosecond + seq + pid + localSalt + macAddr
	src := make([]byte, 8+4+2+localSaltLen+6)

	nowNanosecond := uint64(time.Now().UnixNano())
	src[0] = byte(nowNanosecond >> 56)
	src[1] = byte(nowNanosecond >> 48)
	src[2] = byte(nowNanosecond >> 40)
	src[3] = byte(nowNanosecond >> 32)
	src[4] = byte(nowNanosecond >> 24)
	src[5] = byte(nowNanosecond >> 16)
	src[6] = byte(nowNanosecond >> 8)
	src[7] = byte(nowNanosecond)

	src[8] = byte(seq >> 24)
	src[9] = byte(seq >> 16)
	src[10] = byte(seq >> 8)
	src[11] = byte(seq)

	src[12] = byte(pid >> 8)
	src[13] = byte(pid)

	copy(src[14:], localSalt)
	copy(src[14+localSaltLen:], macAddr)

	hashSum := md5.Sum(src)
	return hashSum[:]
}

// The returned bytes has not been hex encoded, is raw bytes.
func NewRandom() []byte {
	seq := atomic.AddUint32(&randomClockSequence, 1)
	return commonRandom(localRandomSalt, seq)
}

// The returned bytes has been hex encoded.
func NewToken() []byte {
	seq := atomic.AddUint32(&tokenClockSequence, 1)
	token := commonRandom(localTokenSalt, seq)
	ret := make([]byte, hex.EncodedLen(len(token)))
	hex.Encode(ret, token)
	return ret
}

// The returned bytes have been hex encoded.
func NewSessionID() []byte {
	// 32bits unixtime + 48bits mac + 16bits pid + 24bits clockSequence + 40bits md5 sum
	ret := make([]byte, 20)

	// 写入 32bits unixtime
	timenow := time.Now()
	nowSecond := uint32(timenow.Unix())
	ret[0] = byte(nowSecond >> 24)
	ret[1] = byte(nowSecond >> 16)
	ret[2] = byte(nowSecond >> 8)
	ret[3] = byte(nowSecond)

	// 写入 48bits mac
	copy(ret[4:], macAddr)

	// 写入 16bits pid
	ret[10] = byte(pid >> 8)
	ret[11] = byte(pid)

	// 写入 clockSequence
	seq := atomic.AddUint32(&sessionClockSequence, 1)
	ret[12] = byte(seq >> 16)
	ret[13] = byte(seq >> 8)
	ret[14] = byte(seq)

	// 写入 32bits hash sum

	// nowNanosecond + seq + pid + localSalt + macAddr
	salt := make([]byte, 8+4+2+localSaltLen+6)

	nowNanosecond := uint64(timenow.UnixNano())
	salt[0] = byte(nowNanosecond >> 56)
	salt[1] = byte(nowNanosecond >> 48)
	salt[2] = byte(nowNanosecond >> 40)
	salt[3] = byte(nowNanosecond >> 32)
	salt[4] = byte(nowNanosecond >> 24)
	salt[5] = byte(nowNanosecond >> 16)
	salt[6] = byte(nowNanosecond >> 8)
	salt[7] = byte(nowNanosecond)

	salt[8] = byte(seq >> 24)
	salt[9] = byte(seq >> 16)
	salt[10] = byte(seq >> 8)
	salt[11] = byte(seq)

	salt[12] = byte(pid >> 8)
	salt[13] = byte(pid)

	copy(salt[14:], localSessionSalt)
	copy(salt[14+localSaltLen:], macAddr)

	hashSum := md5.Sum(salt)
	ret[15] = hashSum[0]
	ret[16] = hashSum[1]
	ret[17] = hashSum[2]
	ret[18] = hashSum[3]
	ret[19] = hashSum[4]

	hexRet := make([]byte, hex.EncodedLen(len(ret)))
	hex.Encode(hexRet, ret)
	return hexRet
}
