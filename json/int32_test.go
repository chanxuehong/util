package json

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
)

func TestInt32_MarshalJSON(t *testing.T) {
	for _, v := range []Int32{0, 1234, math.MinInt32, math.MaxInt32} {
		data1, err1 := json.Marshal(v)

		data2, err2 := json.Marshal(int32(v))
		data2 = []byte(`"` + string(data2) + `"`)

		if !bytes.Equal(data1, data2) || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
			t.Errorf("failed, value: %v, have(%s, %v), want(%s, %v)", v, data1, err1, data2, err2)
			return
		}
		t.Log(string(data1), err1)
	}
}

func TestInt32_UnmarshalJSON(t *testing.T) {
	// 不带引号
	{
		type T1 struct {
			X Int32 `json:"x"`
		}
		type T2 struct {
			X int32 `json:"x"`
		}
		tests := [][]byte{
			[]byte(`{"x":null}`),
			[]byte(`{"x":0}`),
			[]byte(`{"x":1234}`),
			[]byte(`{"x":-2147483648}`),
			[]byte(`{"x":-2147483649}`),
			[]byte(`{"x":2147483647}`),
			[]byte(`{"x":2147483648}`),
			[]byte(`{"x":-1234}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if int32(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
	// 带引号
	{
		type T1 struct {
			X Int32 `json:"x"`
		}
		type T2 struct {
			X int32 `json:"x,string"`
		}
		tests := [][]byte{
			[]byte(`{"x":"0"}`),
			[]byte(`{"x":"1234"}`),
			[]byte(`{"x":"-2147483648"}`),
			[]byte(`{"x":"-2147483649"}`),
			[]byte(`{"x":"2147483647"}`),
			[]byte(`{"x":"2147483648"}`),
			[]byte(`{"x":"-1234"}`),
			[]byte(`{"x":""}`),
			[]byte(`{"x":"abc"}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if int32(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
}
