package util

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

// 解析xml, 返回第一级子节点的键值集合, 如果第一级子节点包含有子节点, 则跳过.
func ParseXMLToMap(xmlReader io.Reader) (m map[string]string, err error) {
	if xmlReader == nil {
		err = errors.New("nil xmlReader")
		return
	}

	d := xml.NewDecoder(xmlReader)
	m = make(map[string]string)

	var (
		tk    xml.Token    // 当前节点的 xml.Token
		depth int          // 当前节点的深度
		key   string       // 当前"第一级"子节点的 key
		value bytes.Buffer // 当前"第一级"子节点的 value
	)
	for {
		tk, err = d.Token()
		if err != nil {
			if err != io.EOF {
				return
			}
			err = nil
			return
		}

		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 1:
			case 2: // 第一级子节点
				key = v.Name.Local
				value.Reset()
			case 3:
				if err = d.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" 暗示了当前第一级子节点包含子节点
			default:
				panic("incorrect algorithm")
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				value.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = value.String()
			}
			depth--
		}
	}
}

// 格式化 map[string]string 为 xml 格式, 根节点名字为 xml.
//  NOTE: 该函数假定 m map[string]string 里的 key 都是合法的 xml 字符串, 不包含需要转义的字符!
func FormatMapToXML(xmlWriter io.Writer, m map[string]string) (err error) {
	if xmlWriter == nil {
		return errors.New("nil xmlWriter")
	}

	if _, err = io.WriteString(xmlWriter, "<xml>"); err != nil {
		return
	}

	for k, v := range m {
		if _, err = io.WriteString(xmlWriter, "<"+k+">"); err != nil {
			return
		}
		if err = xml.EscapeText(xmlWriter, []byte(v)); err != nil {
			return
		}
		if _, err = io.WriteString(xmlWriter, "</"+k+">"); err != nil {
			return
		}
	}

	if _, err = io.WriteString(xmlWriter, "</xml>"); err != nil {
		return
	}
	return
}
