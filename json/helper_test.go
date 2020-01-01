package json

import (
	"reflect"
	"testing"
)

func TestToStdIntSlice(t *testing.T) {
	tests := []struct {
		src []Int
		dst []int
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Int{},
			dst: []int{},
		},
		{
			src: []Int{1},
			dst: []int{1},
		},
		{
			src: []Int{1, 2, -3, -4},
			dst: []int{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := ToStdIntSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdIntSlice(t *testing.T) {
	tests := []struct {
		src []int
		dst []Int
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []int{},
			dst: []Int{},
		},
		{
			src: []int{1},
			dst: []Int{1},
		},
		{
			src: []int{1, 2, -3, -4},
			dst: []Int{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := FromStdIntSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdUintSlice(t *testing.T) {
	tests := []struct {
		src []Uint
		dst []uint
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Uint{},
			dst: []uint{},
		},
		{
			src: []Uint{1},
			dst: []uint{1},
		},
		{
			src: []Uint{1, 2, 3, 4},
			dst: []uint{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := ToStdUintSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdUintSlice(t *testing.T) {
	tests := []struct {
		src []uint
		dst []Uint
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []uint{},
			dst: []Uint{},
		},
		{
			src: []uint{1},
			dst: []Uint{1},
		},
		{
			src: []uint{1, 2, 3, 4},
			dst: []Uint{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := FromStdUintSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdInt32Slice(t *testing.T) {
	tests := []struct {
		src []Int32
		dst []int32
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Int32{},
			dst: []int32{},
		},
		{
			src: []Int32{1},
			dst: []int32{1},
		},
		{
			src: []Int32{1, 2, -3, -4},
			dst: []int32{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := ToStdInt32Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdInt32Slice(t *testing.T) {
	tests := []struct {
		src []int32
		dst []Int32
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []int32{},
			dst: []Int32{},
		},
		{
			src: []int32{1},
			dst: []Int32{1},
		},
		{
			src: []int32{1, 2, -3, -4},
			dst: []Int32{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := FromStdInt32Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdInt64Slice(t *testing.T) {
	tests := []struct {
		src []Int64
		dst []int64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Int64{},
			dst: []int64{},
		},
		{
			src: []Int64{1},
			dst: []int64{1},
		},
		{
			src: []Int64{1, 2, -3, -4},
			dst: []int64{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := ToStdInt64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdInt64Slice(t *testing.T) {
	tests := []struct {
		src []int64
		dst []Int64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []int64{},
			dst: []Int64{},
		},
		{
			src: []int64{1},
			dst: []Int64{1},
		},
		{
			src: []int64{1, 2, -3, -4},
			dst: []Int64{1, 2, -3, -4},
		},
	}
	for _, v := range tests {
		dst := FromStdInt64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdUint32Slice(t *testing.T) {
	tests := []struct {
		src []Uint32
		dst []uint32
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Uint32{},
			dst: []uint32{},
		},
		{
			src: []Uint32{1},
			dst: []uint32{1},
		},
		{
			src: []Uint32{1, 2, 3, 4},
			dst: []uint32{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := ToStdUint32Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdUint32Slice(t *testing.T) {
	tests := []struct {
		src []uint32
		dst []Uint32
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []uint32{},
			dst: []Uint32{},
		},
		{
			src: []uint32{1},
			dst: []Uint32{1},
		},
		{
			src: []uint32{1, 2, 3, 4},
			dst: []Uint32{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := FromStdUint32Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdUint64Slice(t *testing.T) {
	tests := []struct {
		src []Uint64
		dst []uint64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Uint64{},
			dst: []uint64{},
		},
		{
			src: []Uint64{1},
			dst: []uint64{1},
		},
		{
			src: []Uint64{1, 2, 3, 4},
			dst: []uint64{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := ToStdUint64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdUint64Slice(t *testing.T) {
	tests := []struct {
		src []uint64
		dst []Uint64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []uint64{},
			dst: []Uint64{},
		},
		{
			src: []uint64{1},
			dst: []Uint64{1},
		},
		{
			src: []uint64{1, 2, 3, 4},
			dst: []Uint64{1, 2, 3, 4},
		},
	}
	for _, v := range tests {
		dst := FromStdUint64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdFloat64Slice(t *testing.T) {
	tests := []struct {
		src []Float64
		dst []float64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Float64{},
			dst: []float64{},
		},
		{
			src: []Float64{1},
			dst: []float64{1},
		},
		{
			src: []Float64{1, 2.5, -3, -4.12345},
			dst: []float64{1, 2.5, -3, -4.12345},
		},
	}
	for _, v := range tests {
		dst := ToStdFloat64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdFloat64Slice(t *testing.T) {
	tests := []struct {
		src []float64
		dst []Float64
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []float64{},
			dst: []Float64{},
		},
		{
			src: []float64{1},
			dst: []Float64{1},
		},
		{
			src: []float64{1, 2.5, -3, -4.12345},
			dst: []Float64{1, 2.5, -3, -4.12345},
		},
	}
	for _, v := range tests {
		dst := FromStdFloat64Slice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestToStdBoolSlice(t *testing.T) {
	tests := []struct {
		src []Bool
		dst []bool
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []Bool{},
			dst: []bool{},
		},
		{
			src: []Bool{true},
			dst: []bool{true},
		},
		{
			src: []Bool{true, false, true},
			dst: []bool{true, false, true},
		},
	}
	for _, v := range tests {
		dst := ToStdBoolSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}

func TestFromStdBoolSlice(t *testing.T) {
	tests := []struct {
		src []bool
		dst []Bool
	}{
		{
			src: nil,
			dst: nil,
		},
		{
			src: []bool{},
			dst: []Bool{},
		},
		{
			src: []bool{true},
			dst: []Bool{true},
		},
		{
			src: []bool{true, false, true},
			dst: []Bool{true, false, true},
		},
	}
	for _, v := range tests {
		dst := FromStdBoolSlice(v.src)
		if !reflect.DeepEqual(dst, v.dst) {
			t.Errorf("failed, want equal, have:%+v, want:%+v", dst, v.dst)
			return
		}
	}
}
