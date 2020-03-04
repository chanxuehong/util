package json

import (
	"bytes"
	"encoding/json"
)

// Int32 是 int32 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Int32 int32

var _ json.Marshaler = Int32(0)

// MarshalJSON 实现了 json.Marshaler
func (x Int32) MarshalJSON() ([]byte, error) {
	return marshalInt(int64(x))
}

var _ json.Unmarshaler = (*Int32)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Int32) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, jsonNullLiteral) {
		// no-op
		return nil
	}
	n, err := unmarshalInt(data, "Int32", 32)
	if err != nil {
		return err
	}
	*x = Int32(n)
	return nil
}
