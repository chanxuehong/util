package random

import (
	"crypto/sha1"
	"encoding/base64"
	"sync"
	"time"
)

var (
	idMutex         sync.Mutex
	idLastTimestamp int64
)

// 获取一个不重复的 id, 136 年内基本不会重复.
//  NOTE:
//  返回的是原始数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符,
//  特别的, id 的 url_base64 编码不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
func NewId() (id [12]byte) {
	// 32bits unixtime + 24bits mac hashsum + 16bits pid + 24bits clock sequence

	// 写入 32bits unixtime, 这样跨度 136 年不会重复
	timestamp := time.Now().Unix()
	id[0] = byte(timestamp >> 24)
	id[1] = byte(timestamp >> 16)
	id[2] = byte(timestamp >> 8)
	id[3] = byte(timestamp)

	// 写入 24bits macSHA1HashSum
	copy(id[4:], macSHA1HashSum[:3])

	// 写入 16bits pid
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)

	// 写入 24bit clock sequence, 这样 1 秒内 16777216 个操作都不会重复
	var seq uint32

	idMutex.Lock()
	if timestamp <= idLastTimestamp {
		idClockSequence++
	}
	seq = idClockSequence
	idLastTimestamp = timestamp
	idMutex.Unlock()

	id[9] = byte(seq >> 16)
	id[10] = byte(seq >> 8)
	id[11] = byte(seq)

	return
}

// 返回 time.Time 的 unix 时间, 单位是 100ns.
func unix100ns(t time.Time) uint64 {
	return uint64(t.Unix())*1e7 + uint64(t.Nanosecond())/100
}

var (
	sessionIdMutex         sync.Mutex
	sessionIdLastTimestamp uint64
)

// 获取 sessionid.
// 325 天内基本不会重复(实际上更大的跨度也很难重复), 对于 session 而言这个跨度基本满足了.
//  NOTE:
//  返回的结果是 32 字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
func NewSessionId() (id []byte) {
	timestamp := unix100ns(time.Now())

	// 48bits unix100nano + 48bits mac + 16bits pid + 16bits clock sequence + 64bits SHA1 sum
	var idx [24]byte

	// 写入 48bits unix100nano; 写入低 48 bit, 这样跨度 325 天不会重复
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

	// 写入 16bit clock sequence, 这样 100 纳秒内 65536 个操作都不会重复
	var seq uint64

	sessionIdMutex.Lock()
	if timestamp <= sessionIdLastTimestamp {
		sessionIdClockSequence++
	}
	seq = sessionIdClockSequence
	sessionIdLastTimestamp = timestamp
	sessionIdMutex.Unlock()

	idx[14] = byte(seq >> 8)
	idx[15] = byte(seq)

	// 写入 64bits hashsum, 让 sessionid 猜测的难度增加, 一定程度也能提高唯一性的概率;
	// 特别是 timestamp 轮回325天后出现 timestamp 的低48位 + seq 的低16位和以前的某个时刻刚好相等,
	// 但是这个时候 timestamp 和 seq 和那个时候的不一定相等, localSessionIdSalt 更难相等,
	// 所以后面的 hashsum 就很大可能不相等(SHA1 的碰撞概率很低), 这样还是能保证唯一性!

	var src [8 + 8 + localSaltLen]byte // timestamp + seq + localSessionIdSalt

	src[0] = byte(timestamp >> 56)
	src[1] = byte(timestamp >> 48)
	src[2] = byte(timestamp>>40) ^ localRandomSalt[4]
	src[3] = byte(timestamp>>32) ^ localRandomSalt[5]
	src[4] = byte(timestamp>>24) ^ localRandomSalt[6]
	src[5] = byte(timestamp>>16) ^ localRandomSalt[7]
	src[6] = byte(timestamp>>8) ^ localTokenSalt[4]
	src[7] = byte(timestamp) ^ localTokenSalt[5]

	src[8] = byte(seq >> 56)
	src[9] = byte(seq >> 48)
	src[10] = byte(seq >> 40)
	src[11] = byte(seq >> 32)
	src[12] = byte(seq >> 24)
	src[13] = byte(seq >> 16)
	src[14] = byte(seq>>8) ^ localTokenSalt[6]
	src[15] = byte(seq) ^ localTokenSalt[7]

	copy(src[16:], localSessionIdSalt)

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
