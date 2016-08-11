package json

import (
	"errors"
	"fmt"
	"strconv"
)

type Float64 float64

func (x Float64) MarshalJSON() (data []byte, err error) {
	data = make([]byte, 0, 24+2)
	data = append(data, '"')
	data = strconv.AppendFloat(data, float64(x), 'g', -1, 64)
	data = append(data, '"')
	return
}

func (x *Float64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("json: cannot unmarshal empty string into Go value of type Float64")
	}
	if data[0] != '"' {
		n, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
		}
		*x = Float64(n)
		return nil
	}
	data2 := data[1:]
	if len(data2) == 0 {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
	}
	data2 = data2[:len(data2)-1]
	if len(data2) == 0 {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
	}
	n, err := strconv.ParseFloat(string(data2), 64)
	if err != nil {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
	}
	*x = Float64(n)
	return nil
}
