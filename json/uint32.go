package json

import (
	"encoding/json"
)

// Uint32 是 uint32 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Uint32 uint32

var _ json.Marshaler = Uint32(0)

// MarshalJSON 实现了 json.Marshaler
func (x Uint32) MarshalJSON() ([]byte, error) {
	return marshalUint(uint64(x))
}

var _ json.Unmarshaler = (*Uint32)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Uint32) UnmarshalJSON(data []byte) error {
	n, err := unmarshalUint(data, "Uint32", 32)
	if err != nil {
		return err
	}
	*x = Uint32(n)
	return nil
}
