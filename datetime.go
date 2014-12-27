package util

import (
	"time"
)

var BeijingLocation = time.FixedZone("UTC+8", 8*60*60)

const (
	secondsBeijingOffset = 8 * 60 * 60
	secondsPerDay        = 24 * 60 * 60
)

// time.Time 转换为距 1970-01-01 的天数(UTC+8)
func TimeToBeijingUnixDay(t time.Time) int64 {
	return (t.Unix() + secondsBeijingOffset) / secondsPerDay
}

// 距 1970-01-01 的天数(UTC+8) 转换为 time.Time, BeijingLocation
func BeijingUnixDayToTime(n int64) time.Time {
	return time.Unix(n*secondsPerDay-secondsBeijingOffset, 0).In(BeijingLocation)
}

// unixtime 转换为 time.Time, BeijingLocation
func UnixToBeijingLocationTime(n int64) time.Time {
	return time.Unix(n, 0).In(BeijingLocation)
}
