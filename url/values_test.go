package url

import (
	"testing"
)

type EncodeQueryTest struct {
	m        Values
	expected string
}

var encodeQueryTests = []EncodeQueryTest{
	{
		m:        nil,
		expected: "",
	},
	{
		m:        Values{"q": {"puppies"}, "oe": {"utf8"}},
		expected: "oe=utf8&q=puppies",
	},
	{
		m:        Values{"q": {"dogs", "&", "7"}},
		expected: "q=dogs&q=%26&q=7",
	},
	{
		m: Values{
			"a": {"a1", "a2", "a3"},
			"b": {"b1", "b2", "b3"},
			"c": {"c1", "c2", "c3"},
		},
		expected: "a=a1&a=a2&a=a3&b=b1&b=b2&b=b3&c=c1&c=c2&c=c3",
	},
	{
		m:        Values{"q": {"puppies"}, "oe": {"utf8"}, "new": {"hello world"}},
		expected: "new=hello%20world&oe=utf8&q=puppies",
	},
}

func TestEncodeQuery(t *testing.T) {
	for _, tt := range encodeQueryTests {
		if q := tt.m.Encode(); q != tt.expected {
			t.Errorf(`EncodeQuery(%+v) = %q, want %q`, tt.m, q, tt.expected)
		}
	}
}
