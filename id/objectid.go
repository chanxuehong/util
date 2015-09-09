package id

import (
	"errors"
	"sync"
	"time"

	"github.com/chanxuehong/util/random"
)

// 返回从 2010-01-01 00:00:00 +0000 UTC 到 t time.Time 所经过的毫秒的个数
func idtimestamp(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/1000000 - 1262304000000
}

// tillNextMillis spin wait till next millisecond.
func tillNextMillis(lastTimestamp int64) int64 {
	timestamp := idtimestamp(time.Now())
	for timestamp <= lastTimestamp {
		timestamp = idtimestamp(time.Now())
	}
	return timestamp
}

const objectIdSequenceMask = 0x3fff // 14bits

var (
	objectIdMutex    sync.Mutex
	objectIdSequence uint32
	// 最近的时间戳的第一个 sequence.
	// 对于同一个时间戳, 如果 objectIdSequence 再次等于 objectIdFirstSequence,
	// 表示达到了上限了, 需要等到下一个时间戳了.
	objectIdFirstSequence uint32
	objectIdLastTimestamp int64
)

func init() {
	objectIdSequence = random.NewRandomUint32() & objectIdSequenceMask
}

// 获取一个不重复的 id (每毫秒可以产生 16384 个 id, 和 mongodb 的 objectid 算法类似, 不完全一致).
//  NOTE:
//  1. 从 2010-01-01 00:00:00 +0000 UTC 到 2149-05-15 07:35:11.104 +0000 UTC 时间段内生成的 id 是升序且不重复的.
//  2. 返回的 id 是原始字节数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符,
//  3. 特别的, id 的 url_base64 编码不包含等号(=), 只有 1-9,a-z,A-Z,-,_ 字符.
//  4. 这个 id 适合在自己的系统内部用, 如果想要给外部用最要用 uuid.ver1.
func NewObjectId() (id [12]byte, err error) {
	timestamp := idtimestamp(time.Now())

	objectIdMutex.Lock() // Lock
	switch {
	case timestamp > objectIdLastTimestamp:
		objectIdFirstSequence = objectIdSequence
	case timestamp == objectIdLastTimestamp:
		objectIdSequence = (objectIdSequence + 1) & objectIdSequenceMask
		if objectIdSequence == objectIdFirstSequence {
			timestamp = tillNextMillis(timestamp)
		}
	default:
		objectIdMutex.Unlock() // Unlock
		err = errors.New("Clock moved backwards")
		return
	}
	objectIdLastTimestamp = timestamp
	sequence := objectIdSequence
	objectIdMutex.Unlock() // Unlock

	// 42bits timestamp + 6bits higher sequence + 24bits mac hashsum + 16bits pid + 8bits lower sequence

	// 写入 42bits timestamp, 这样跨度 139 年不会重复
	id[0] = byte(timestamp >> 34)
	id[1] = byte(timestamp >> 26)
	id[2] = byte(timestamp >> 18)
	id[3] = byte(timestamp >> 10)
	id[4] = byte(timestamp >> 2)
	id[5] = byte(timestamp << 6)

	// 写入 sequence 的高 6 位
	id[5] |= byte(sequence>>8) & 0x3f

	// 写入 24bits mac hashsum
	copy(id[6:], macHashSum[:3])

	// 写入 16bits pid
	id[9] = byte(pid >> 8)
	id[10] = byte(pid)

	// 写入 sequence 的低 8 位,
	// 加上前面的高 6 位, 这样就是 1 毫秒内可以产生 16384 个 id
	id[11] = byte(sequence)
	return
}
