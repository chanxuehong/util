package id

import (
	"errors"
	"sync"
	"time"

	"github.com/chanxuehong/util/random"
)

const snowflakeSequenceMask = 0xfff // 12bits

var (
	snowflakeWorkerId        int64 = -1 // 10bits, [0, 1024)
	snowflakeSetWorkerIdOnce sync.Once

	snowflakeMutex    sync.Mutex
	snowflakeSequence uint32 = random.NewRandomUint32() & snowflakeSequenceMask
	// 最近的时间戳的第一个 sequence.
	// 对于同一个时间戳, 如果 snowflakeSequence 再次等于 snowflakeFirstSequence,
	// 表示达到了上限了, 需要等到下一个时间戳了.
	snowflakeFirstSequence uint32
	snowflakeLastTimestamp int64
)

// 设置 snowflake worker Id, [0, 1024).
func SetSnowflakeWorkerId(workerId int) (err error) {
	if workerId < 0 || workerId > 1023 {
		return errors.New("worker Id can't be greater than 1023 or less than 0")
	}
	snowflakeSetWorkerIdOnce.Do(func() {
		snowflakeWorkerId = int64(workerId)
	})
	return
}

// 获取一个不重复的 id (每毫秒可以产生 4096 个 id, snowflake, 纪元不一样).
//  NOTE:
//  1. 从 2010-01-01 00:00:00 +0000 UTC 到 2079-09-07 15:47:35.552 +0000 UTC 时间段内生成的 id 是升序且不重复的.
//  2. 这个 id 适合在自己的系统内部用, 否则最好用 uuid.ver1.
func NewSnowflakeId() (id int64, err error) {
	if snowflakeWorkerId == -1 {
		err = errors.New("WorkerId has not been assigned")
		return
	}
	timestamp := idtimestamp(time.Now())

	snowflakeMutex.Lock() // Lock
	switch {
	case timestamp > snowflakeLastTimestamp:
		snowflakeFirstSequence = snowflakeSequence
	case timestamp == snowflakeLastTimestamp:
		snowflakeSequence = (snowflakeSequence + 1) & snowflakeSequenceMask
		if snowflakeSequence == snowflakeFirstSequence {
			timestamp = tillNextMillis(timestamp)
		}
	default:
		snowflakeMutex.Unlock() // Unlock
		err = errors.New("Clock moved backwards")
		return
	}
	snowflakeLastTimestamp = timestamp
	sequence := snowflakeSequence
	snowflakeMutex.Unlock() // Unlock

	// 0(1bit) + 41bits timestamp + 10bits worker id + 12bits sequence
	id = timestamp<<22 | snowflakeWorkerId<<12 | int64(sequence)
	id &= 0x7fffffffffffffff
	return
}
