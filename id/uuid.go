package id

import (
	"sync/atomic"
	"time"

	"github.com/chanxuehong/util/random"
)

const (
	unixToUUID = 122192928000000000 // 从 1582-10-15 00:00:00 +0000 UTC 到 1970-01-01 00:00:00 +0000 UTC 的 100ns 的个数
)

// 返回 uuid 的时间戳, 从 1582-10-15 00:00:00 +0000 UTC 到 t time.Time 的 100ns 的个数.
func uuid100ns(t time.Time) uint64 {
	return unix100ns(t) + unixToUUID
}

// 返回 unix 的时间戳, 从 1970-01-01 00:00:00 +0000 UTC 到 t time.Time 的 100ns 的个数.
func unix100ns(t time.Time) uint64 {
	return uint64(t.Unix())*10000000 + uint64(t.Nanosecond())/100
}

var uuidSequence uint32 = random.NewRandomUint32()

// 返回 uuid, ver1.
//  NOTE: 返回的是原始字节数组, 不是可显示字符, 可以通过 hex, url_base64 等转换为可显示字符.
func NewUUIDV1() (uuid [16]byte) {
	timestamp := uuid100ns(time.Now())

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
	seq := atomic.AddUint32(&uuidSequence, 1)
	uuid[8] = byte(seq>>8) & 0x3F
	uuid[9] = byte(seq)

	// set variant
	uuid[8] |= 0x80

	// set node, 48bits
	copy(uuid[10:], realMAC[:])
	return
}
