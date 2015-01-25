package util

import (
	"time"
)

var BeijingLocation = time.FixedZone("UTC+8", 8*60*60)

const (
	secondsPerDay = 24 * 60 * 60

	utcToBeijing = 8 * 60 * 60
	beijingToUtc = -utcToBeijing
)

// UTC unixtime 转换为北京时间的 unixtime.
func UTCUnixToBeijingUnix(n int64) int64 {
	return n + utcToBeijing
}

// time.Time 转换为北京时间的 unixtime.
func TimeToBeijingUnix(t time.Time) int64 {
	return UTCUnixToBeijingUnix(t.Unix())
}

// 北京时间的 unixtime 转换为 UTC unixtime.
func BeijingUnixToUTCUnix(n int64) int64 {
	return n + beijingToUtc
}

// 北京时间的 unixtime 转换为 time.Time, BeijingLocation.
func BeijingUnixToTime(n int64) time.Time {
	return time.Unix(BeijingUnixToUTCUnix(n), 0).In(BeijingLocation)
}

// 北京时间的 unixtime 转换为北京时间距 1970-01-01 的天数.
func BeijingUnixToBeijingUnixDay(n int64) int64 {
	return n / secondsPerDay
}

// UTC unixtime 转换为北京时间距 1970-01-01 的天数.
func UTCUnixToBeijingUnixDay(n int64) int64 {
	return BeijingUnixToBeijingUnixDay(UTCUnixToBeijingUnix(n))
}

// time.Time 转换为北京时间距 1970-01-01 的天数.
func TimeToBeijingUnixDay(t time.Time) int64 {
	return BeijingUnixToBeijingUnixDay(TimeToBeijingUnix(t))
}

// 北京时间的 unixtime 转换为 北京时间的 unixtime.
func BeijingUnixDayToBeijingUnix(n int64) int64 {
	return n * secondsPerDay
}

// 北京时间距 1970-01-01 的天数转换为 UTC unixtime.
func BeijingUnixDayToUTCUnix(n int64) int64 {
	return BeijingUnixToUTCUnix(BeijingUnixDayToBeijingUnix(n))
}

// 北京时间距 1970-01-01 的天数转换为 time.Time, BeijingLocation.
func BeijingUnixDayToTime(n int64) time.Time {
	return BeijingUnixToTime(BeijingUnixDayToBeijingUnix(n))
}
