package time

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Time 扩展了标准 time.Time.
// 实现  database/sql/driver.Valuer, database/sql.Scanner 接口, 数据库字段类型为 int64, unixtime;
// 重新实现了 encoding/json.Marshaler, encoding/json.Unmarshaler 接口, 格式为 "2006-01-02 15:04:05", 北京时间
type Time struct {
	time.Time
}

func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{
		Time: time.Date(year, month, day, hour, min, sec, nsec, loc),
	}
}

func Now() Time {
	return Time{
		Time: time.Now(),
	}
}

func Parse(layout, value string) (t Time, err error) {
	tt, err := time.Parse(layout, value)
	if err != nil {
		return
	}
	t = Time{
		Time: tt,
	}
	return
}

func ParseInLocation(layout, value string, loc *time.Location) (t Time, err error) {
	tt, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return
	}
	t = Time{
		Time: tt,
	}
	return
}

func Unix(sec int64, nsec int64) Time {
	return Time{
		Time: time.Unix(sec, nsec),
	}
}

func (t Time) Value() (value driver.Value, err error) {
	value = t.Unix()
	return
}

func (t *Time) Scan(value interface{}) (err error) {
	var unixtime int64
	if err = convertAssign(&unixtime, value); err != nil {
		return
	}
	t.Time = time.Unix(unixtime, 0)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y > 9999 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.In(BeijingLocation).Format("2006-01-02 15:04:05")), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.ParseInLocation("2006-01-02 15:04:05", string(data), BeijingLocation)
	return
}
