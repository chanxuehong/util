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
