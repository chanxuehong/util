package id

import (
	"crypto/sha1"
	"encoding/base64"
	"sync"
	"time"

	"github.com/chanxuehong/util/random"
)

const (
	sessionIdSaltLen            = 39   // see NewSessionId(), 8+8+39 == 55, sha1 签名性能最好的前提下最大的数据块
	sessionIdSaltUpdateInterval = 3600 // seconds
)

var (
	sessionIdSalt               []byte = make([]byte, sessionIdSaltLen)
	sessionIdMutex              sync.Mutex
	sessionIdSaltLastUpdateTime int64  = time.Now().Unix()
	sessionIdClockSequence      uint64 = random.NewRandomUint64()
)

func init() {
	random.ReadRandomBytes(sessionIdSalt)
}

// 获取 sessionid.
//  NOTE:
//  1. 325 天内基本不会重复(实际上更大的跨度也很难重复), 对于 session 而言这个跨度基本满足了.
//  2. 返回的结果是 32 字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
func NewSessionId() (id []byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()
	timestamp := unix100ns(timeNow)

	sessionIdMutex.Lock() // Lock
	if timeNowUnix >= sessionIdSaltLastUpdateTime+sessionIdSaltUpdateInterval {
		sessionIdSaltLastUpdateTime = timeNowUnix
		random.ReadRandomBytes(sessionIdSalt)
	}
	sessionIdClockSequence++
	clockSequence := sessionIdClockSequence
	sessionIdMutex.Unlock() // Unlock

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
	copy(idx[6:], fakeMAC[:])

	// 写入 16bits pid
	idx[12] = byte(pid >> 8)
	idx[13] = byte(pid)

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
	src[2] = byte(timestamp>>40) ^ sessionIdSalt[0]
	src[3] = byte(timestamp>>32) ^ sessionIdSalt[1]
	src[4] = byte(timestamp>>24) ^ sessionIdSalt[2]
	src[5] = byte(timestamp>>16) ^ sessionIdSalt[3]
	src[6] = byte(timestamp>>8) ^ sessionIdSalt[4]
	src[7] = byte(timestamp) ^ sessionIdSalt[5]

	src[8] = byte(clockSequence >> 56)
	src[9] = byte(clockSequence >> 48)
	src[10] = byte(clockSequence >> 40)
	src[11] = byte(clockSequence >> 32)
	src[12] = byte(clockSequence >> 24)
	src[13] = byte(clockSequence >> 16)
	src[14] = byte(clockSequence>>8) ^ sessionIdSalt[6]
	src[15] = byte(clockSequence) ^ sessionIdSalt[7]

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
