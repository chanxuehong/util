package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseXMLToMap(t *testing.T) {
	var xmlSrc = []string{
		`<xml>
			<a>a</a>
			<b>b</b>
		<xml>`,
		`<xml>
			<a>a</a>
			<b>
				<ba>ba</ba>
			</b>
			<c>c</c>
		<xml>`,
		`<xml>
			<a>a</a>
			<b>
				bchara
				<ba>ba</ba>
			</b>
			<c>c</c>
		<xml>`,
		`<xml>
			<a>a</a>
			<b>
				bchara
				<ba>ba</ba>
				bchara
				<bb>bb</bb>
				bchara
			</b>
			<c>c</c>
		<xml>`,
		`<xml>
			chara
			<a>a</a>
			<b>
				<ba>ba</ba>
				bchara
			</b>
			<c>c</c>
		<xml>`,
	}

	var mapWant = []map[string]string{
		{
			"a": "a",
			"b": "b",
		},
		{
			"a": "a",
			"c": "c",
		},
		{
			"a": "a",
			"c": "c",
		},
		{
			"a": "a",
			"c": "c",
		},
		{
			"a": "a",
			"c": "c",
		},
	}

	for i, src := range xmlSrc {
		m, err := ParseXMLToMap(strings.NewReader(src))
		if err != nil {
			t.Errorf("ParseXMLToMap(%s):\nError: %s\n", src, err.Error())
			continue
		}
		if !reflect.DeepEqual(m, mapWant[i]) {
			t.Errorf("ParseXMLToMap(%s):\nhave %v\nwant %v\n", src, m, mapWant[i])
			continue
		}
	}
}
