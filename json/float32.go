package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Float32 是 float32 的 json 类型, marshal 会变成字符串, unmarshal 可以接受数字或字符串.
type Float32 float32

var _ json.Marshaler = Float32(0)

// MarshalJSON 实现了 json.Marshaler
func (x Float32) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(float32(x))
	if err != nil {
		return nil, err
	}
	var buf [64]byte
	buf[0] = '"'
	result := buf[:1]
	result = append(result, data...)
	result = append(result, '"')
	return result, nil
}

var _ json.Unmarshaler = (*Float32)(nil)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Float32) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return errors.New("json: cannot unmarshal empty string into Go value of type Float32")
	}
	if bytes.Equal(data, jsonNullLiteral) {
		*x = Float32(0)
		return nil
	}
	if data[0] != '"' {
		n, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return fmt.Errorf("json: cannot unmarshal string %q into Go value of type Float32", data)
		}
		*x = Float32(n)
		return nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return fmt.Errorf("json: cannot unmarshal string %q into Go value of type Float32", data)
	}
	n, err := strconv.ParseFloat(string(data[1:maxIndex]), 32)
	if err != nil {
		return fmt.Errorf("json: cannot unmarshal string %q into Go value of type Float32", data)
	}
	*x = Float32(n)
	return nil
}
