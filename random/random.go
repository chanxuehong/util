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
	nowNanosecond := time.Now().UnixNano()

	// 56bits unix nanosecond time + 48bits mac + 24bits clock sequence + 32bits md5 sum
	ret := make([]byte, 20)

	// 写入 56bits unix nanosecond time; 以 100 纳秒为单位, 写入低 56 bit
	nowNanosecondX := nowNanosecond / 100
	ret[0] = byte(nowNanosecondX >> 48)
	ret[1] = byte(nowNanosecondX >> 40)
	ret[2] = byte(nowNanosecondX >> 32)
	ret[3] = byte(nowNanosecondX >> 24)
	ret[4] = byte(nowNanosecondX >> 16)
	ret[5] = byte(nowNanosecondX >> 8)
	ret[6] = byte(nowNanosecondX)

	// 写入 48bits mac
	copy(ret[7:], macAddr)

	// 写入 24bit clock sequence
	seq := atomic.AddUint32(&sessionClockSequence, 1)
	ret[13] = byte(seq >> 16)
	ret[14] = byte(seq >> 8)
	ret[15] = byte(seq)

	// 写入 32bits hash sum
	salt := make([]byte, 8+4+localSaltLen+6) // nowNanosecond^pid + seq + localSessionSalt + macAddr
	salt[0] = byte(nowNanosecond>>56) ^ byte(pid>>8)
	salt[1] = byte(nowNanosecond>>48) ^ byte(pid)
	salt[2] = byte(nowNanosecond>>40) ^ byte(pid>>8)
	salt[3] = byte(nowNanosecond>>32) ^ byte(pid)
	salt[4] = byte(nowNanosecond>>24) ^ byte(pid>>8)
	salt[5] = byte(nowNanosecond>>16) ^ byte(pid)
	salt[6] = byte(nowNanosecond>>8) ^ byte(pid>>8)
	salt[7] = byte(nowNanosecond) ^ byte(pid)

	salt[8] = byte(seq >> 24)
	salt[9] = byte(seq >> 16)
	salt[10] = byte(seq >> 8)
	salt[11] = byte(seq)

	copy(salt[12:], localSessionSalt)
	copy(salt[12+localSaltLen:], macAddr)

	hashSum := md5.Sum(salt)
	ret[16] = hashSum[0]
	ret[17] = hashSum[1]
	ret[18] = hashSum[2]
	ret[19] = hashSum[3]

	hexRet := make([]byte, hex.EncodedLen(len(ret)))
	hex.Encode(hexRet, ret)
	return hexRet
}
