package url

import (
	"net/url"
	"strings"
	"testing"
)

func TestQueryEscape(t *testing.T) {
	buf := make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		buf = append(buf, byte(i))
	}
	buf = append(buf, "你好世界"...)
	str := string(buf)

	str1 := url.QueryEscape(str)
	str2 := QueryEscape(str)

	if str1 == str2 {
		t.Errorf("want not equal")
		return
	}

	str1 = strings.ReplaceAll(str1, "+", "%20")
	if str1 != str2 {
		t.Errorf("want equal, but have:%q, want:%q", str2, str1)
		return
	}
}

func TestQueryUnescape(t *testing.T) {
	buf := make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		buf = append(buf, byte(i))
	}
	buf = append(buf, "你好世界"...)
	str := string(buf)

	{
		str2, err := QueryUnescape(QueryEscape(str))
		if err != nil {
			t.Error(err.Error())
			return
		}
		if str != str2 {
			t.Errorf("want equal, but have:%q, want:%q", str2, str)
			return
		}
	}

	{
		str2, err := QueryUnescape(url.QueryEscape(str))
		if err != nil {
			t.Error(err.Error())
			return
		}
		if str != str2 {
			t.Errorf("want equal, but have:%q, want:%q", str2, str)
			return
		}
	}

	{
		str2, err := url.QueryUnescape(QueryEscape(str))
		if err != nil {
			t.Error(err.Error())
			return
		}
		if str != str2 {
			t.Errorf("want equal, but have:%q, want:%q", str2, str)
			return
		}
	}

	{
		str2, err := url.QueryUnescape(url.QueryEscape(str))
		if err != nil {
			t.Error(err.Error())
			return
		}
		if str != str2 {
			t.Errorf("want equal, but have:%q, want:%q", str2, str)
			return
		}
	}
}
