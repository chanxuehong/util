package json

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

func TestUint_MarshalJSON(t *testing.T) {
	tests := []Uint{0, 1234}
	if strconv.IntSize == 32 {
		tests = append(tests, math.MaxUint32)
	} else {
		tests = append(tests, math.MaxUint64)
	}
	for _, v := range tests {
		data1, err1 := json.Marshal(v)

		data2, err2 := json.Marshal(uint(v))
		data2 = []byte(`"` + string(data2) + `"`)

		if !bytes.Equal(data1, data2) || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
			t.Errorf("failed, value: %v, have(%s, %v), want(%s, %v)", v, data1, err1, data2, err2)
			return
		}
		t.Log(string(data1), err1)
	}
}

func TestUint_UnmarshalJSON(t *testing.T) {
	var maxUintLiteral string
	var maxUintOverflowLiteral string
	if strconv.IntSize == 32 {
		maxUintLiteral = "4294967295"
		maxUintOverflowLiteral = "4294967296"
	} else {
		maxUintLiteral = "18446744073709551615"
		maxUintOverflowLiteral = "18446744073709551616"
	}

	// 不带引号
	{
		type T1 struct {
			X Uint `json:"x"`
		}
		type T2 struct {
			X uint `json:"x"`
		}
		tests := [][]byte{
			[]byte(`{"x":null}`),
			[]byte(`{"x":0}`),
			[]byte(`{"x":1234}`),
			[]byte(`{"x":` + maxUintLiteral + `}`),
			[]byte(`{"x":` + maxUintOverflowLiteral + `}`),
			[]byte(`{"x":-1234}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if uint(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
	// 带引号
	{
		type T1 struct {
			X Uint `json:"x"`
		}
		type T2 struct {
			X uint `json:"x,string"`
		}
		tests := [][]byte{
			[]byte(`{"x":"0"}`),
			[]byte(`{"x":"1234"}`),
			[]byte(`{"x":"` + maxUintLiteral + `"}`),
			[]byte(`{"x":"` + maxUintOverflowLiteral + `"}`),
			[]byte(`{"x":"-1234"}`),
			[]byte(`{"x":""}`),
			[]byte(`{"x":"abc"}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if uint(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%d, %v), want(%d, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
}
