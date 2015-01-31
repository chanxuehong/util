package random

import (
	"sync/atomic"
	"time"
)

const (
	unixToUUID = 122192928000000000 // 从 1582-10-15T00:00:00 到 1970-01-01T00:00:00 的 100ns 的个数
)

// 返回 uuid 的时间戳, 从 1582-10-15T00:00:00 到 time.Time 的 100ns 的个数.
func uuid100ns(t time.Time) uint64 {
	return unix100ns(t) + unixToUUID
}

// 返回 uuid, version == 1.
//  NOTE: 返回的是原始数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符.
func NewUUIDV1() (u [16]byte) {
	timestamp := uuid100ns(time.Now())

	// set timestamp, 60bits
	u[0] = byte(timestamp >> 24)
	u[1] = byte(timestamp >> 16)
	u[2] = byte(timestamp >> 8)
	u[3] = byte(timestamp)

	u[4] = byte(timestamp >> 40)
	u[5] = byte(timestamp >> 32)

	u[6] = byte(timestamp>>56) & 0x0F
	u[7] = byte(timestamp >> 48)

	// set version, 4bits
	u[6] |= 0x10

	// set clock sequence, 14bits
	seq := atomic.AddUint32(&uuidClockSequence, 1)
	u[8] = byte(seq>>8) & 0x3F
	u[9] = byte(seq)

	// set variant
	u[8] |= 0x80

	// set node, 48bits
	copy(u[10:], mac[:])
	return
}
