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

	var key string          // 当前"第一级"子节点的 key
	var buffer bytes.Buffer // 当前"第一级"子节点的 value
	var depth int           // 当前节点的深度

	m = make(map[string]string)
	for {
		var tk xml.Token
		tk, err = d.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}

		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 1: // do nothing
			case 2:
				key = v.Name.Local
				buffer.Reset()
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
				buffer.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = buffer.String()
			}
			depth--
		}
	}
}

// 格式化 map[string]string 为 xml 格式, 根节点名字为 xml.
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
