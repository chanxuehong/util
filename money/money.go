package money

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Money 表示金钱, 单位为分.
type Money int64

var _ fmt.Stringer = Money(0)

var (
	_ encoding.TextMarshaler   = Money(0)
	_ encoding.TextUnmarshaler = (*Money)(nil)
)

var (
	_ json.Marshaler   = Money(0)
	_ json.Unmarshaler = (*Money)(nil)
)

var (
	_ xml.Marshaler   = Money(0)
	_ xml.Unmarshaler = (*Money)(nil)
)

// String 将 Money 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money) String() string {
	text, _ := m.MarshalText()
	return string(text)
}

// MarshalText 将 Money 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money) MarshalText() (text []byte, err error) {
	text, _ = m.MarshalJSON()
	return text[1 : len(text)-1], nil
}

// MarshalJSON 将 Money 编码成 "xxxx.yz" 这样以 '元' 为单位的字符串.
func (m Money) MarshalJSON() (text []byte, err error) {
	switch {
	case m > 0:
		str := strconv.FormatInt(int64(m), 10)
		switch len(str) {
		case 1: // x --> "0.0x"
			bs := make([]byte, 0, 6)
			bs = append(bs, `"0.0`...)
			bs = append(bs, str[0])
			bs = append(bs, '"')
			return bs, nil
		case 2: // xy --> "0.xy"
			bs := make([]byte, 0, 6)
			bs = append(bs, `"0.`...)
			bs = append(bs, str...)
			bs = append(bs, '"')
			return bs, nil
		default: // len(str) >= 3
			if strings.HasSuffix(str, "00") { // xxxx00 --> "xxxx"
				bs := make([]byte, 0, len(str))
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '"')
				return bs, nil
			}
			// xxxxyz --> "xxxx.yz"
			bs := make([]byte, 0, len(str)+3)
			bs = append(bs, '"')
			bs = append(bs, str[:len(str)-2]...)
			bs = append(bs, '.')
			bs = append(bs, str[len(str)-2:]...)
			bs = append(bs, '"')
			return bs, nil
		}
	case m == 0:
		return []byte{'"', '0', '"'}, nil
	default: // n < 0
		str := strconv.FormatInt(int64(m), 10)
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
			if strings.HasSuffix(str, "00") { // -xxxx00 --> "-xxxx"
				bs := make([]byte, 0, len(str))
				bs = append(bs, '"')
				bs = append(bs, str[:len(str)-2]...)
				bs = append(bs, '"')
				return bs, nil
			}
			// -xxxxyz --> "-xxxx.yz"
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

// MarshalXML 将 Money 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return
	}
	text, _ := m.MarshalText()
	if err = e.EncodeToken(xml.CharData(text)); err != nil {
		return
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalText 将 xxxx.yz 这样以 '元' 为单位的字符串解码到 Money 中.
func (m *Money) UnmarshalText(text []byte) (err error) {
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
	*m = Money(n)
	return nil
}

// UnmarshalJSON 将 "xxxx.yz" 这样以 '元' 为单位的字符串解码到 Money 中.
func (m *Money) UnmarshalJSON(data []byte) (err error) {
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[0] != '"' || data[maxIndex] != '"' {
		return fmt.Errorf("invalid Money JSON text: %s", data)
	}
	return m.UnmarshalText(data[1:maxIndex])
}

// UnmarshalXML 将 xxxx.yz 这样以 '元' 为单位的字符串解码到 Money 中.
func (m *Money) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var text []byte
	if err = d.DecodeElement(&text, &start); err != nil {
		return
	}
	return m.UnmarshalText(text)
}
