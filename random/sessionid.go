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

// 获取 sessionid, 理论上 325天 内不会重复, 对于 session 而言这个跨度基本满足了.
//  NOTE: 返回的结果已经经过 base64 url 编码
func NewSessionId() (id []byte) {
	timenow := time.Now()

	// 48bits unix100nano + 48bits mac + 16bits pid + 16bits clock sequence + 64bits SHA-1 sum
	idx := make([]byte, 24)

	// 写入 48bits unix100nano; 写入低 48 bit, 这样跨度 325天 不会重复.
	timestamp := unix100nano(timenow)
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
	seq := atomic.AddUint32(&sessionClockSequence, 1)
	idx[14] = byte(seq >> 8)
	idx[15] = byte(seq)

	// 写入 64bits hash sum, 让 sessionid 猜测的难度增加

	hashSrc := make([]byte, 8+4+localSaltLen) // timestamp + seq + localSessionSalt

	// 因为 idx 开头暴露了 timestamp, 所以这里要混淆下
	hashSrc[0] = byte(timestamp>>56) ^ localRandomSalt[0]
	hashSrc[1] = byte(timestamp>>48) ^ localRandomSalt[1]
	hashSrc[2] = byte(timestamp>>40) ^ localRandomSalt[2]
	hashSrc[3] = byte(timestamp>>32) ^ localRandomSalt[3]
	hashSrc[4] = byte(timestamp>>24) ^ localTokenSalt[0]
	hashSrc[5] = byte(timestamp>>16) ^ localTokenSalt[1]
	hashSrc[6] = byte(timestamp>>8) ^ localTokenSalt[2]
	hashSrc[7] = byte(timestamp) ^ localTokenSalt[3]

	hashSrc[8] = byte(seq >> 24)
	hashSrc[9] = byte(seq >> 16)
	hashSrc[10] = byte(seq >> 8)
	hashSrc[11] = byte(seq)

	copy(hashSrc[12:], localSessionSalt)

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
