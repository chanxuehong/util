package json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Bool 是 bool 的 json 类型, unmarshal 可以接受 boolean 或 字符串.
type Bool bool

var _ json.Unmarshaler = (*Bool)(nil)

var (
	jsonNullLiteral        = []byte(`null`)
	jsonTrueLiteral        = []byte(`true`)
	jsonFalseLiteral       = []byte(`false`)
	jsonTrueQuotedLiteral  = []byte(`"true"`)
	jsonFalseQuotedLiteral = []byte(`"false"`)
)

// UnmarshalJSON 实现了 json.Unmarshaler
func (x *Bool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	switch len(data) {
	case 4:
		if bytes.Equal(data, jsonTrueLiteral) {
			*x = Bool(true)
			return nil
		}
		if bytes.Equal(data, jsonNullLiteral) {
			*x = Bool(false)
			return nil
		}
	case 5:
		if bytes.Equal(data, jsonFalseLiteral) {
			*x = Bool(false)
			return nil
		}
	case 6:
		if bytes.Equal(data, jsonTrueQuotedLiteral) {
			*x = Bool(true)
			return nil
		}
	case 7:
		if bytes.Equal(data, jsonFalseQuotedLiteral) {
			*x = Bool(false)
			return nil
		}
	}
	return fmt.Errorf("json: cannot unmarshal string %q into Go value of type Bool", data)
}
