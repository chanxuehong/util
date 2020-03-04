package json

import (
	"bytes"
	"encoding/json"
)

// Uint64 是 uint64 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Uint64 uint64

var _ json.Marshaler = Uint64(0)

// MarshalJSON 实现了 json.Marshaler
func (x Uint64) MarshalJSON() ([]byte, error) {
	return marshalUint(uint64(x))
}

var _ json.Unmarshaler = (*Uint64)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Uint64) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, jsonNullLiteral) {
		// no-op
		return nil
	}
	n, err := unmarshalUint(data, "Uint64", 64)
	if err != nil {
		return err
	}
	*x = Uint64(n)
	return nil
}
