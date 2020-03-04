package json

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Int 是 int 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Int int

var _ json.Marshaler = Int(0)

// MarshalJSON 实现了 json.Marshaler
func (x Int) MarshalJSON() ([]byte, error) {
	return marshalInt(int64(x))
}

var _ json.Unmarshaler = (*Int)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Int) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, jsonNullLiteral) {
		// no-op
		return nil
	}
	n, err := unmarshalInt(data, "Int", strconv.IntSize)
	if err != nil {
		return err
	}
	*x = Int(n)
	return nil
}
