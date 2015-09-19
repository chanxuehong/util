package timex

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	time.Time
	Valid bool // Valid is true if Time is not NULL
}

// unixtime
func (nt NullTime) Value() (value driver.Value, err error) {
	if !nt.Valid {
		return
	}
	value = nt.Unix()
	return
}

// unixtime
func (nt *NullTime) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}
	nt.Valid = true
	var unixtime int64
	if err = convertAssign(&unixtime, value); err != nil {
		return
	}
	nt.Time = time.Unix(unixtime, 0)
	return
}
