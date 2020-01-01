package json

import (
	"encoding/json"
	"testing"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	// 不带引号
	{
		type T1 struct {
			X Bool `json:"x"`
		}
		type T2 struct {
			X bool `json:"x"`
		}
		tests := [][]byte{
			[]byte(`{"x":null}`),
			[]byte(`{"x":false}`),
			[]byte(`{"x":true}`),
			[]byte(`{"x":False}`),
			[]byte(`{"x":True}`),
			[]byte(`{"x":FALSE}`),
			[]byte(`{"x":TRUE}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if bool(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%t, %v), want(%t, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
	// 带引号
	{
		type T1 struct {
			X Bool `json:"x"`
		}
		type T2 struct {
			X bool `json:"x,string"`
		}
		tests := [][]byte{
			[]byte(`{"x":"false"}`),
			[]byte(`{"x":"true"}`),
			[]byte(`{"x":"False"}`),
			[]byte(`{"x":"True"}`),
			[]byte(`{"x":"FALSE"}`),
			[]byte(`{"x":"TRUE"}`),
			[]byte(`{"x":""}`),
			[]byte(`{"x":"abc"}`),
		}
		for _, data := range tests {
			var v1 T1
			err1 := json.Unmarshal(data, &v1)

			var v2 T2
			err2 := json.Unmarshal(data, &v2)

			if bool(v1.X) != v2.X || (err1 != nil && err2 == nil) || (err1 == nil && err2 != nil) {
				t.Errorf("failed, data: %s, have(%t, %v), want(%t, %v)", data, v1.X, err1, v2.X, err2)
				return
			}
			t.Log(v1.X, err1, "------", v2.X, err2)
		}
	}
}
