package json

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Uint 是 uint 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Uint uint

var _ json.Marshaler = Uint(0)

// MarshalJSON 实现了 json.Marshaler
func (x Uint) MarshalJSON() ([]byte, error) {
	return marshalUint(uint64(x))
}

var _ json.Unmarshaler = (*Uint)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Uint) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, jsonNullLiteral) {
		// no-op
		return nil
	}
	n, err := unmarshalUint(data, "Uint", strconv.IntSize)
	if err != nil {
		return err
	}
	*x = Uint(n)
	return nil
}
