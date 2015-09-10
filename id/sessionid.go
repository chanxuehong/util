package id

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/chanxuehong/util/random"
)

const (
	sessionIdSaltLen            = 39   // see NewSessionId(), 8+8+39 == 55, sha1 签名性能最好的前提下最大的数据块
	sessionIdSaltUpdateInterval = 3600 // seconds

	sessionIdSequenceMask = 0xffff // 16bits
)

var (
	sessionIdSalt []byte = make([]byte, sessionIdSaltLen)

	sessionIdMutex                   sync.Mutex
	sessionIdSaltLastUpdateTimestamp int64

	sessionIdSequence uint64 = random.NewRandomUint64()
	// 最近的时间戳的第一个 sequence.
	// 对于同一个时间戳, 如果 sessionIdSequence 再次等于 sessionIdFirstSequence,
	// 表示达到了上限了, 需要等到下一个时间戳了.
	sessionIdFirstSequence uint64
	sessionIdLastTimestamp int64
)

// sessionIdTillNext100ns spin wait till next 100 nanosecond.
func sessionIdTillNext100ns(lastTimestamp int64) int64 {
	timestamp := unix100ns(time.Now())
	for timestamp <= lastTimestamp {
		timestamp = unix100ns(time.Now())
	}
	return timestamp
}

// 获取 sessionid.
//  NOTE:
//  1. 325 天内基本不会重复(实际上更大的跨度也很难重复), 对于 session 而言这个跨度基本满足了.
//  2. 返回的结果是 32 字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
func NewSessionId() (id []byte, err error) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()
	timestamp := unix100ns(timeNow)

	isSaltUpdated := false

	sessionIdMutex.Lock() // Lock
	switch {
	case timestamp > sessionIdLastTimestamp:
		sessionIdSequence++
		sessionIdFirstSequence = sessionIdSequence & sessionIdSequenceMask
	case timestamp == sessionIdLastTimestamp:
		sessionIdSequence++
		if sessionIdSequence&sessionIdSequenceMask == sessionIdFirstSequence {
			timestamp = sessionIdTillNext100ns(timestamp)
		}
	default:
		sessionIdMutex.Unlock() // Unlock
		err = errors.New("Clock moved backwards")
		return
	}
	sessionIdLastTimestamp = timestamp
	sequence := sessionIdSequence

	if timeNowUnix >= sessionIdSaltLastUpdateTimestamp+sessionIdSaltUpdateInterval {
		sessionIdSaltLastUpdateTimestamp = timeNowUnix
		random.ReadRandomBytes(sessionIdSalt)
		isSaltUpdated = true
	}
	sessionIdMutex.Unlock() // Unlock

	// 48bits unix100ns + 48bits mac + 16bits pid + 16bits sequence + 64bits SHA1 sum
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

	// 写入 16bit sequence, 这样 100 纳秒内 65536 个操作都不会重复
	idx[14] = byte(sequence >> 8)
	idx[15] = byte(sequence)

	if isSaltUpdated {
		copy(idx[16:], sessionIdSalt)
	} else {
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

		src[8] = byte(sequence >> 56)
		src[9] = byte(sequence >> 48)
		src[10] = byte(sequence >> 40)
		src[11] = byte(sequence >> 32)
		src[12] = byte(sequence >> 24)
		src[13] = byte(sequence >> 16)
		src[14] = byte(sequence>>8) ^ sessionIdSalt[6]
		src[15] = byte(sequence) ^ sessionIdSalt[7]

		copy(src[16:], sessionIdSalt)

		hashsum := sha1.Sum(src[:])
		copy(idx[16:], hashsum[:])
	}
	id = make([]byte, 32)
	base64.URLEncoding.Encode(id, idx[:])
	return
}
