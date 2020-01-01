package json

import (
	"encoding/json"
)

// Int64 是 int64 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Int64 int64

var _ json.Marshaler = Int64(0)

// MarshalJSON 实现了 json.Marshaler
func (x Int64) MarshalJSON() ([]byte, error) {
	return marshalInt(int64(x))
}

var _ json.Unmarshaler = (*Int64)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Int64) UnmarshalJSON(data []byte) error {
	n, err := unmarshalInt(data, "Int64", 64)
	if err != nil {
		return err
	}
	*x = Int64(n)
	return nil
}
