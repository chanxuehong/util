package random

import (
	"crypto/sha1"
	"encoding/base64"
	"sync/atomic"
	"time"
)

// 返回参数 t time.Time 的 unix 时间, 单位是 100 纳秒
func unix100nano(t time.Time) int64 {
	return t.Unix()*1e7 + int64(t.Nanosecond()/100)
}

// 获取一个不重复的 id, 58455年 内基本不会重复.
//  NOTE: 返回的结果是 24字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_
func NewId() (id []byte) {
	timestamp := unix100nano(time.Now())

	// 64bits unix100nano + 48bits mac + 16bits pid + 16bits clock sequence
	idx := make([]byte, 18)

	// 写入 64bits unix100nano, 这样跨度 58455年 不会重复
	idx[0] = byte(timestamp >> 56)
	idx[1] = byte(timestamp >> 48)
	idx[2] = byte(timestamp >> 40)
	idx[3] = byte(timestamp >> 32)
	idx[4] = byte(timestamp >> 24)
	idx[5] = byte(timestamp >> 16)
	idx[6] = byte(timestamp >> 8)
	idx[7] = byte(timestamp)

	// 写入 48bits mac
	copy(idx[8:], macAddr)

	// 写入 16bits pid
	idx[14] = byte(pid >> 8)
	idx[15] = byte(pid)

	// 写入 16bit clock sequence, 这样 100 纳秒内 65536 个操作都不会重复
	seq := atomic.AddUint32(&idClockSequence, 1)
	idx[16] = byte(seq >> 8)
	idx[17] = byte(seq)

	id = make([]byte, 24)
	base64.URLEncoding.Encode(id, idx)
	return
}

// 获取 sessionid.
// 325天 内基本不会重复(实际上更大的跨度也很难重复), 对于 session 而言这个跨度基本满足了.
//  NOTE: 返回的结果是 32字节的 url base64 编码, 不包含等号(=), 只有 1-9,a-z,A-Z,-,_
func NewSessionId() (id []byte) {
	timestamp := unix100nano(time.Now())

	// 48bits unix100nano + 48bits mac + 16bits pid + 16bits clock sequence + 64bits SHA-1 sum
	idx := make([]byte, 24)

	// 写入 48bits unix100nano; 写入低 48 bit, 这样跨度 325天 不会重复
	idx[0] = byte(timestamp >> 40)
	idx[1] = byte(timestamp >> 32)
	idx[2] = byte(timestamp >> 24)
	idx[3] = byte(timestamp >> 16)
	idx[4] = byte(timestamp >> 8)
	idx[5] = byte(timestamp)

	// 写入 48bits mac
	copy(idx[6:], macAddr)

	// 写入 16bits pid
	idx[12] = byte(pid >> 8)
	idx[13] = byte(pid)

	// 写入 16bit clock sequence, 这样 100 纳秒内 65536 个操作都不会重复
	seq := atomic.AddUint64(&sessionIdClockSequence, 1)
	idx[14] = byte(seq >> 8)
	idx[15] = byte(seq)

	// 写入 64bits hash sum, 让 sessionid 猜测的难度增加;
	// 一定程度也能提高唯一性的概率, 特别是 timestamp 轮回(325天)后出现 timestamp+seq的低16位
	// 和以前的某个时刻刚好相等, 但是这个时候 seq 和那个时候的 seq 不一定相等, 所以后面的
	// hashsum 就很大可能不相等(SHA-1 的碰撞概率很低), 这样还是能保证唯一性!

	hashSrc := make([]byte, 8+8+localSaltLen) // timestamp + seq + localSessionSalt

	// 因为 idx 开头暴露了 timestamp, 所以这里要混淆下
	hashSrc[0] = byte(timestamp>>56) ^ localRandomSalt[0]
	hashSrc[1] = byte(timestamp>>48) ^ localRandomSalt[1]
	hashSrc[2] = byte(timestamp>>40) ^ localRandomSalt[2]
	hashSrc[3] = byte(timestamp>>32) ^ localRandomSalt[3]
	hashSrc[4] = byte(timestamp>>24) ^ localTokenSalt[0]
	hashSrc[5] = byte(timestamp>>16) ^ localTokenSalt[1]
	hashSrc[6] = byte(timestamp>>8) ^ localTokenSalt[2]
	hashSrc[7] = byte(timestamp) ^ localTokenSalt[3]

	// seq 整个 64bits 都写入进去
	hashSrc[8] = byte(seq >> 56)
	hashSrc[9] = byte(seq >> 48)
	hashSrc[10] = byte(seq >> 40)
	hashSrc[11] = byte(seq >> 32)
	hashSrc[12] = byte(seq >> 24)
	hashSrc[13] = byte(seq >> 16)
	hashSrc[14] = byte(seq >> 8)
	hashSrc[15] = byte(seq)

	copy(hashSrc[16:], localSessionSalt)

	hashSum := sha1.Sum(hashSrc)

	idx[16] = hashSum[12]
	idx[17] = hashSum[13]
	idx[18] = hashSum[14]
	idx[19] = hashSum[15]
	idx[20] = hashSum[16]
	idx[21] = hashSum[17]
	idx[22] = hashSum[18]
	idx[23] = hashSum[19]

	id = make([]byte, 32)
	base64.URLEncoding.Encode(id, idx)
	return
}
