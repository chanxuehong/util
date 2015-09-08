package id

import (
	"sync/atomic"
	"time"

	"github.com/chanxuehong/util/random"
)

var objectIdClockSequence uint32 = random.NewRandomUint32()

// 获取一个不重复的 id, 136 年内基本不会重复(和 mongodb 的 objectid 算法一致).
//  NOTE:
//  1. 从 1970-01-01 00:00:00 +0000 UTC 到 2106-02-07 06:28:15 +0000 UTC 时间段内生成的 id 是升序的
//  2. 返回的 id 是原始字节数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符,
//  3. 特别的, id 的 url_base64 编码不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
//  4. 这个 id 适合在自己的系统内部用, 如果想要给外部用最要用 uuid.ver1.
func NewObjectId() (id [12]byte) {
	// 32bits unixtime + 24bits mac hashsum + 16bits pid + 24bits clock sequence

	// 写入 32bits unixtime, 这样跨度 136 年不会重复
	timestamp := time.Now().Unix()
	id[0] = byte(timestamp >> 24)
	id[1] = byte(timestamp >> 16)
	id[2] = byte(timestamp >> 8)
	id[3] = byte(timestamp)

	// 写入 24bits mac hashsum
	copy(id[4:], macSHA1HashSum[:3])

	// 写入 16bits pid
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)

	// 写入 24bit clock sequence, 这样 1 秒内 16777216 个操作都不会重复
	seq := atomic.AddUint32(&objectIdClockSequence, 1)
	id[9] = byte(seq >> 16)
	id[10] = byte(seq >> 8)
	id[11] = byte(seq)
	return
}
