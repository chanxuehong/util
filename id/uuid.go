package id

import (
	"errors"
	"sync"
	"time"

	"github.com/chanxuehong/util/random"
)

const unixToUUID = 122192928000000000 // 从 1582-10-15 00:00:00 +0000 UTC 到 1970-01-01 00:00:00 +0000 UTC 的 100ns 的个数

// 返回 uuid 的时间戳, 从 1582-10-15 00:00:00 +0000 UTC 到 t time.Time 的 100ns 的个数.
func uuid100ns(t time.Time) int64 {
	return unix100ns(t) + unixToUUID
}

// 返回 unix 的时间戳, 从 1970-01-01 00:00:00 +0000 UTC 到 t time.Time 的 100ns 的个数.
func unix100ns(t time.Time) int64 {
	return t.Unix()*10000000 + int64(t.Nanosecond())/100
}

// uuidTillNext100ns spin wait till next 100 nanosecond.
func uuidTillNext100ns(lastTimestamp int64) int64 {
	timestamp := uuid100ns(time.Now())
	for timestamp <= lastTimestamp {
		timestamp = uuid100ns(time.Now())
	}
	return timestamp
}

const uuidSequenceMask = 0x3fff // 14bits

var (
	uuidMutex    sync.Mutex
	uuidSequence uint32
	// 最近的时间戳的第一个 sequence.
	// 对于同一个时间戳, 如果 uuidSequence 再次等于 uuidFirstSequence,
	// 表示达到了上限了, 需要等到下一个时间戳了.
	uuidFirstSequence uint32
	uuidLastTimestamp int64
)

func init() {
	uuidSequence = random.NewRandomUint32() & uuidSequenceMask
}

// 返回 uuid, ver1.
//  NOTE: 返回的是原始字节数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符.
func NewUUIDV1() (uuid [16]byte, err error) {
	timestamp := uuid100ns(time.Now())

	uuidMutex.Lock() // Lock
	switch {
	case timestamp > uuidLastTimestamp:
		uuidFirstSequence = uuidSequence
	case timestamp == uuidLastTimestamp:
		uuidSequence = (uuidSequence + 1) & uuidSequenceMask
		if uuidSequence == uuidFirstSequence {
			timestamp = uuidTillNext100ns(timestamp)
		}
	default:
		uuidMutex.Unlock() // Unlock
		err = errors.New("Clock moved backwards")
		return
	}
	uuidLastTimestamp = timestamp
	sequence := uuidSequence
	uuidMutex.Unlock() // Unlock

	// set timestamp, 60bits
	uuid[0] = byte(timestamp >> 24)
	uuid[1] = byte(timestamp >> 16)
	uuid[2] = byte(timestamp >> 8)
	uuid[3] = byte(timestamp)

	uuid[4] = byte(timestamp >> 40)
	uuid[5] = byte(timestamp >> 32)

	uuid[6] = byte(timestamp>>56) & 0x0F
	uuid[7] = byte(timestamp >> 48)

	// set version, 4bits
	uuid[6] |= 0x10

	// set sequence, 14bits
	uuid[8] = byte(sequence>>8) & 0x3F
	uuid[9] = byte(sequence)

	// set variant
	uuid[8] |= 0x80

	// set node, 48bits
	copy(uuid[10:], realMAC[:])
	return
}
