package random

import (
	"crypto/sha1"
	"encoding/hex"
	"sync/atomic"
	"time"
)

// 返回参数 t time.Time 的 unix 时间, 单位是 100 纳秒
func unix100nano(t time.Time) int64 {
	return t.Unix()*1e7 + int64(t.Nanosecond()/100)
}

// 获取 sessionid, 理论上 325天 内不会重复, 对于 session 而言这个跨度基本满足了.
//  NOTE: 返回的结果已经经过 hex 编码
func NewSessionId() []byte {
	timenow := time.Now()

	// 48bits unix100nano + 48bits mac + 16bits pid + 16bits clock sequence + 64bits SHA-1 sum
	ret := make([]byte, 24)

	// 写入 48bits unix100nano; 写入低 48 bit, 这样跨度 325天 不会重复.
	timestamp := unix100nano(timenow)
	ret[0] = byte(timestamp >> 40)
	ret[1] = byte(timestamp >> 32)
	ret[2] = byte(timestamp >> 24)
	ret[3] = byte(timestamp >> 16)
	ret[4] = byte(timestamp >> 8)
	ret[5] = byte(timestamp)

	// 写入 48bits mac
	copy(ret[6:], macAddr)

	// 写入 16bits pid
	ret[12] = byte(pid >> 8)
	ret[13] = byte(pid)

	// 写入 16bit clock sequence, 这样 100 纳秒内 65536 个操作都不会重复
	seq := atomic.AddUint32(&sessionClockSequence, 1)
	ret[14] = byte(seq >> 8)
	ret[15] = byte(seq)

	// 写入 64bits hash sum, 让 sessionid 猜测的难度增加

	hashSrc := make([]byte, 8+4+localSaltLen) // timestamp + seq + localSessionSalt

	// 因为 ret 开头暴露了 timestamp, 所以这里要混淆下
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

	ret[16] = hashSum[12]
	ret[17] = hashSum[13]
	ret[18] = hashSum[14]
	ret[19] = hashSum[15]
	ret[20] = hashSum[16]
	ret[21] = hashSum[17]
	ret[22] = hashSum[18]
	ret[23] = hashSum[19]

	hexRet := make([]byte, hex.EncodedLen(len(ret)))
	hex.Encode(hexRet, ret)
	return hexRet
}
