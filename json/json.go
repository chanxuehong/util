package json

import (
	"fmt"
	"strconv"
)

func marshalInt(n int64) ([]byte, error) {
	data := make([]byte, 0, 20+2)
	data = append(data, '"')
	data = strconv.AppendInt(data, n, 10)
	data = append(data, '"')
	return data, nil
}

func marshalUint(n uint64) ([]byte, error) {
	data := make([]byte, 0, 20+2)
	data = append(data, '"')
	data = strconv.AppendUint(data, n, 10)
	data = append(data, '"')
	return data, nil
}

func unmarshalInt(data []byte, typeName string, bitSize int) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("json: cannot unmarshal empty string into Go value of type %s", typeName)
	}
	if data[0] != '"' {
		n, err := strconv.ParseInt(string(data), 10, bitSize)
		if err != nil {
			return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
		}
		return n, nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
	}
	n, err := strconv.ParseInt(string(data[1:maxIndex]), 10, bitSize)
	if err != nil {
		return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
	}
	return n, nil
}

func unmarshalUint(data []byte, typeName string, bitSize int) (uint64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("json: cannot unmarshal empty string into Go value of type %s", typeName)
	}
	if data[0] != '"' {
		n, err := strconv.ParseUint(string(data), 10, bitSize)
		if err != nil {
			return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
		}
		return n, nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
	}
	n, err := strconv.ParseUint(string(data[1:maxIndex]), 10, bitSize)
	if err != nil {
		return 0, fmt.Errorf("json: cannot unmarshal string %q into Go value of type %s", data, typeName)
	}
	return n, nil
}
