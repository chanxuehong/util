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

// 100 纳秒为单位的 unix time 时间戳
func unix100NanoTimestamp(t time.Time) int64 {
	return t.Unix()*1e7 + int64(t.Nanosecond()/100)
}

// The returned bytes have been hex encoded.
func NewSessionID() []byte {
	timenow := time.Now()

	// 56bits unix 100*nanosecond time + 48bits mac + 16bits pid + 24bits clock sequence + 48bits md5 sum
	ret := make([]byte, 24)

	// 写入 56bits unix 100*nanosecond time;
	// 以 100 纳秒为单位, 写入低 56 bit, 这样跨度 228 年不会重复.
	timestamp := unix100NanoTimestamp(timenow)
	ret[0] = byte(timestamp >> 48)
	ret[1] = byte(timestamp >> 40)
	ret[2] = byte(timestamp >> 32)
	ret[3] = byte(timestamp >> 24)
	ret[4] = byte(timestamp >> 16)
	ret[5] = byte(timestamp >> 8)
	ret[6] = byte(timestamp)

	// 写入 48bits mac
	copy(ret[7:], macAddr)

	// 写入 16bits pid
	ret[13] = byte(pid >> 8)
	ret[14] = byte(pid)

	// 写入 24bit clock sequence, 这样 100 纳秒内 16777216 个操作都不会重复
	seq := atomic.AddUint32(&sessionClockSequence, 1)
	ret[15] = byte(seq >> 16)
	ret[16] = byte(seq >> 8)
	ret[17] = byte(seq)

	// 写入 48bits hash sum
	salt := make([]byte, 8+4+2+localSaltLen+6) // nowNanosecond + seq + pid + localSessionSalt + macAddr
	// 因为 ret 开头暴露了 nowNanosecond, 所以这里要混淆下
	nowNanosecond := timenow.UnixNano()
	salt[0] = byte(nowNanosecond>>56) ^ localRandomSalt[4]
	salt[1] = byte(nowNanosecond>>48) ^ localRandomSalt[5]
	salt[2] = byte(nowNanosecond>>40) ^ localRandomSalt[6]
	salt[3] = byte(nowNanosecond>>32) ^ localRandomSalt[7]
	salt[4] = byte(nowNanosecond>>24) ^ localTokenSalt[4]
	salt[5] = byte(nowNanosecond>>16) ^ localTokenSalt[5]
	salt[6] = byte(nowNanosecond>>8) ^ localTokenSalt[6]
	salt[7] = byte(nowNanosecond) ^ localTokenSalt[7]

	salt[8] = byte(seq >> 24)
	salt[9] = byte(seq >> 16)
	salt[10] = byte(seq >> 8)
	salt[11] = byte(seq)

	salt[12] = byte(pid >> 8)
	salt[13] = byte(pid)

	copy(salt[14:], localSessionSalt)
	copy(salt[14+localSaltLen:], macAddr)

	hashSum := md5.Sum(salt)
	ret[18] = hashSum[0]
	ret[19] = hashSum[1]
	ret[20] = hashSum[2]
	ret[21] = hashSum[3]
	ret[22] = hashSum[4]
	ret[23] = hashSum[5]

	hexRet := make([]byte, hex.EncodedLen(len(ret)))
	hex.Encode(hexRet, ret)
	return hexRet
}
