package json

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

func TestInt_MarshalJSON(t *testing.T) {
	tests := []Int{0, 1234}
	if strconv.IntSize == 32 {
		tests = append(tests, math.MinInt32, math.MaxInt32)
	} else {
		tests = append(tests, math.MinInt64, math.MaxInt64)
	}
	for _, v := range tests {
		data1, err1 := json.Marshal(v)

		data2, err2 := json.Marshal(int(v))
		data2 = []byte(`"` + string(data2) + `"`)

		if !bytes.Equal(data1, data2) || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
			t.Errorf("failed, value: %v, have(%s, %v), want(%s, %v)", v, data1, err1, data2, err2)
			return
		}
		t.Log(string(data1), err1)
	}
}

func TestInt_UnmarshalJSON(t *testing.T) {
	var minIntLiteral string
	var minIntOverflowLiteral string
	var maxIntLiteral string
	var maxIntOverflowLiteral string
	if strconv.IntSize == 32 {
		minIntLiteral = "-2147483648"
		minIntOverflowLiteral = "-2147483649"
		maxIntLiteral = "2147483647"
		maxIntOverflowLiteral = "2147483648"
	} else {
		minIntLiteral = "-9223372036854775808"
		minIntOverflowLiteral = "-9223372036854775809"
		maxIntLiteral = "9223372036854775807"
		maxIntOverflowLiteral = "9223372036854775808"
	}

	// 不带引号
	{
		type T1 struct {
			X Int `json:"x"`
		}
		type T2 struct {
			X int `json:"x"`
		}
		tests := [][]byte{
			[]byte(`{"x":null}`),
			[]byte(`{"x":0}`),
			[]byte(`{"x":1234}`),
			[]byte(`{"x":` + minIntLiteral + `}`),
			[]byte(`{"x":` + minIntOverflowLiteral + `}`),
			[]byte(`{"x":` + maxIntLiteral + `}`),
			[]byte(`{"x":` + maxIntOverflowLiteral + `}`),
			[]byte(`{"x":-1234}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if int(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
	// 带引号
	{
		type T1 struct {
			X Int `json:"x"`
		}
		type T2 struct {
			X int `json:"x,string"`
		}
		tests := [][]byte{
			[]byte(`{"x":"0"}`),
			[]byte(`{"x":"1234"}`),
			[]byte(`{"x":"` + minIntLiteral + `"}`),
			[]byte(`{"x":"` + minIntOverflowLiteral + `"}`),
			[]byte(`{"x":"` + maxIntLiteral + `"}`),
			[]byte(`{"x":"` + maxIntOverflowLiteral + `"}`),
			[]byte(`{"x":"-1234"}`),
			[]byte(`{"x":""}`),
			[]byte(`{"x":"abc"}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if int(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
}
