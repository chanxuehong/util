package money

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Money int64 // 基本单位: 分

var _ fmt.Stringer = Money(0)

var (
	_ json.Marshaler   = Money(0)
	_ json.Unmarshaler = (*Money)(nil)
)

var (
	_ encoding.TextMarshaler   = Money(0)
	_ encoding.TextUnmarshaler = (*Money)(nil)
)

func (p Money) String() string {
	text, _ := p.MarshalJSON()
	return string(text[1 : len(text)-1])
}

// MarshalJSON 将 Money 编码成 "xxxx.yz" 这样以 元 为单位的字符串.
func (p Money) MarshalJSON() ([]byte, error) {
	switch n := int64(p); {
	case n > 0:
		str := strconv.FormatInt(n, 10)
		switch len(str) {
		case 1: // x --> "0.0x"
			bs := make([]byte, 0, 6)
			bs = append(bs, `"0.0`...)
			bs = append(bs, str...)
			bs = append(bs, '"')
			return bs, nil
		case 2: // xy --> "0.xy"
			bs := make([]byte, 0, 6)
			bs = append(bs, `"0.`...)
			bs = append(bs, str...)
			bs = append(bs, '"')
			return bs, nil
		default: // len(str) >= 3
			if strings.HasSuffix(str, "00") { // ...x00 --> "...x"
				bs := make([]byte, 0, len(str))
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '"')
				return bs, nil
			} else { // ...xyz --> "...x.yz"
				bs := make([]byte, 0, len(str)+3)
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '.')
				bs = append(bs, str[len(str)-2:]...)
				bs = append(bs, '"')
				return bs, nil
			}
		}
	case n == 0:
		return []byte{'"', '0', '"'}, nil
	default: // n < 0
		str := strconv.FormatInt(n, 10)
		switch len(str) {
		case 2: // -x --> "-0.0x"
			bs := make([]byte, 0, 7)
			bs = append(bs, `"-0.0`...)
			bs = append(bs, str[1])
			bs = append(bs, '"')
			return bs, nil
		case 3: // -xy --> "-0.xy"
			bs := make([]byte, 0, 7)
			bs = append(bs, `"-0.`...)
			bs = append(bs, str[1:]...)
			bs = append(bs, '"')
			return bs, nil
		default: // len(str) >= 4
			if strings.HasSuffix(str, "00") { // -...x00 --> "-...x"
				bs := make([]byte, 0, len(str))
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '"')
				return bs, nil
			} else { // -...xyz --> "-...x.yz"
				bs := make([]byte, 0, len(str)+3)
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '.')
				bs = append(bs, str[len(str)-2:]...)
				bs = append(bs, '"')
				return bs, nil
			}
		}
	}
}

// MarshalText 将 Money 编码成 xxxx.yz 这样以 元 为单位的字符串.
func (p Money) MarshalText() (text []byte, err error) {
	text, _ = p.MarshalJSON()
	return text[1 : len(text)-1], nil
}

// UnmarshalJSON 将 "xxxx.yz" 这样以 元 为单位的字符串解码到 Money 中.
func (p *Money) UnmarshalJSON(data []byte) error {
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[0] != '"' || data[maxIndex] != '"' {
		return fmt.Errorf("invalid Money JSON text: %s", data)
	}
	return p.UnmarshalText(data[1:maxIndex])
}

// UnmarshalText 将 xxxx.yz 这样以 元 为单位的字符串解码到 Money 中.
func (p *Money) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return fmt.Errorf("invalid Money text: %s", text)
	}
	str := string(text)
	if dotIndex := strings.IndexByte(str, '.'); dotIndex >= 0 {
		dotFront := str[:dotIndex]
		dotBehind := str[dotIndex+1:]
		switch len(dotBehind) {
		default:
			return fmt.Errorf("invalid Money text: %s", text)
		case 0:
			str = dotFront + "00"
		case 1:
			str = dotFront + dotBehind + "0"
		case 2:
			str = dotFront + dotBehind
		}
	} else {
		str += "00"
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Money text: %s", text)
	}
	*p = Money(n)
	return nil
}
