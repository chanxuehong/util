package random

import (
	"crypto/sha1"
	"encoding/base64"
	"sync"
	"time"
)

// 返回 time.Time 的 unix 时间, 单位是 100ns.
func unix100ns(t time.Time) uint64 {
	return uint64(t.Unix())*1e7 + uint64(t.Nanosecond())/100
}

const sessionIdSaltLen = 39

var (
	sessionIdSalt = make([]byte, sessionIdSaltLen)

	sessionIdMutex         sync.Mutex
	sessionIdLastTimestamp uint64
	sessionIdClockSequence uint64
	sessionIdRandomFactor  uint64
)

func init() {
	readRandomBytes(sessionIdSalt)
	sessionIdClockSequence = newRandomUint64()
	sessionIdRandomFactor = newRandomUint64()
}

// 获取 sessionid.
// 325 天内基本不会重复(实际上更大的跨度也很难重复), 对于 session 而言这个跨度基本满足了.
//  NOTE:
//  返回的结果是 32 字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
func NewSessionId() (id []byte) {
	timestamp := unix100ns(time.Now())

	// 48bits unix100ns + 48bits mac + 16bits pid + 16bits clock sequence + 64bits SHA1 sum
	var idx [24]byte

	// 写入 unix100ns 低 48 bit, 这样跨度 325 天不会重复
	idx[0] = byte(timestamp >> 40)
	idx[1] = byte(timestamp >> 32)
	idx[2] = byte(timestamp >> 24)
	idx[3] = byte(timestamp >> 16)
	idx[4] = byte(timestamp >> 8)
	idx[5] = byte(timestamp)

	// 写入 48bits mac
	copy(idx[6:], mac[:])

	// 写入 16bits pid
	idx[12] = byte(pid >> 8)
	idx[13] = byte(pid)

	sessionIdMutex.Lock()
	if timestamp <= sessionIdLastTimestamp {
		sessionIdClockSequence++
	}
	clockSequence := sessionIdClockSequence
	sessionIdRandomFactor++
	randomFactor := sessionIdRandomFactor
	sessionIdLastTimestamp = timestamp
	sessionIdMutex.Unlock()

	// 写入 16bit clock sequence, 这样 100 纳秒内 65536 个操作都不会重复
	idx[14] = byte(clockSequence >> 8)
	idx[15] = byte(clockSequence)

	// 写入 64bits hashsum, 让 sessionid 猜测的难度增加, 一定程度也能提高唯一性的概率;
	// 特别是 timestamp 轮回325天后出现 timestamp 的低48位 + seq 的低16位和以前的某个时刻刚好相等,
	// 但是这个时候 timestamp 和 seq 和那个时候的不一定相等, sessionIdSalt 更难相等,
	// 所以后面的 hashsum 就很大可能不相等(SHA1 的碰撞概率很低), 这样还是能保证唯一性!

	var src [8 + 8 + sessionIdSaltLen]byte // 8+8+39 == 55

	src[0] = byte(timestamp >> 56)
	src[1] = byte(timestamp >> 48)
	src[2] = byte(timestamp>>40) ^ byte(randomFactor>>56)
	src[3] = byte(timestamp>>32) ^ byte(randomFactor>>48)
	src[4] = byte(timestamp>>24) ^ byte(randomFactor>>40)
	src[5] = byte(timestamp>>16) ^ byte(randomFactor>>32)
	src[6] = byte(timestamp>>8) ^ byte(randomFactor>>24)
	src[7] = byte(timestamp) ^ byte(randomFactor>>16)

	src[8] = byte(clockSequence >> 56)
	src[9] = byte(clockSequence >> 48)
	src[10] = byte(clockSequence >> 40)
	src[11] = byte(clockSequence >> 32)
	src[12] = byte(clockSequence >> 24)
	src[13] = byte(clockSequence >> 16)
	src[14] = byte(clockSequence>>8) ^ byte(randomFactor>>8)
	src[15] = byte(clockSequence) ^ byte(randomFactor)

	copy(src[16:], sessionIdSalt)

	hashSum := sha1.Sum(src[:])

	idx[16] = hashSum[0]
	idx[17] = hashSum[1]
	idx[18] = hashSum[2]
	idx[19] = hashSum[3]
	idx[20] = hashSum[4]
	idx[21] = hashSum[5]
	idx[22] = hashSum[6]
	idx[23] = hashSum[7]

	id = make([]byte, 32)
	base64.URLEncoding.Encode(id, idx[:])
	return
}
