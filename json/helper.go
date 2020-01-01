package json

import "unsafe"

// ToStdIntSlice cast []Int to []int
func ToStdIntSlice(s []Int) []int {
	return *((*[]int)(unsafe.Pointer(&s)))
}

// FromStdIntSlice cast []int to []Int
func FromStdIntSlice(s []int) []Int {
	return *((*[]Int)(unsafe.Pointer(&s)))
}

// ToStdUintSlice cast []Uint to []uint
func ToStdUintSlice(s []Uint) []uint {
	return *((*[]uint)(unsafe.Pointer(&s)))
}

// FromStdUintSlice cast []uint to []Uint
func FromStdUintSlice(s []uint) []Uint {
	return *((*[]Uint)(unsafe.Pointer(&s)))
}

// ToStdInt32Slice cast []Int32 to []int32
func ToStdInt32Slice(s []Int32) []int32 {
	return *((*[]int32)(unsafe.Pointer(&s)))
}

// FromStdInt32Slice cast []int32 to []Int32
func FromStdInt32Slice(s []int32) []Int32 {
	return *((*[]Int32)(unsafe.Pointer(&s)))
}

// ToStdInt64Slice cast []Int64 to []int64
func ToStdInt64Slice(s []Int64) []int64 {
	return *((*[]int64)(unsafe.Pointer(&s)))
}

// FromStdInt64Slice cast []int64 to []Int64
func FromStdInt64Slice(s []int64) []Int64 {
	return *((*[]Int64)(unsafe.Pointer(&s)))
}

// ToStdUint32Slice cast []Uint32 to []uint32
func ToStdUint32Slice(s []Uint32) []uint32 {
	return *((*[]uint32)(unsafe.Pointer(&s)))
}

// FromStdUint32Slice cast []uint32 to []Uint32
func FromStdUint32Slice(s []uint32) []Uint32 {
	return *((*[]Uint32)(unsafe.Pointer(&s)))
}

// ToStdUint64Slice cast []Uint64 to []uint64
func ToStdUint64Slice(s []Uint64) []uint64 {
	return *((*[]uint64)(unsafe.Pointer(&s)))
}

// FromStdUint64Slice cast []uint64 to []Uint64
func FromStdUint64Slice(s []uint64) []Uint64 {
	return *((*[]Uint64)(unsafe.Pointer(&s)))
}

// ToStdFloat64Slice cast []Float64 to []float64
func ToStdFloat64Slice(s []Float64) []float64 {
	return *((*[]float64)(unsafe.Pointer(&s)))
}

// FromStdFloat64Slice cast []float64 to []Float64
func FromStdFloat64Slice(s []float64) []Float64 {
	return *((*[]Float64)(unsafe.Pointer(&s)))
}

// ToStdBoolSlice cast []Bool to []bool
func ToStdBoolSlice(s []Bool) []bool {
	return *((*[]bool)(unsafe.Pointer(&s)))
}

// FromStdBoolSlice cast []bool to []Bool
func FromStdBoolSlice(s []bool) []Bool {
	return *((*[]Bool)(unsafe.Pointer(&s)))
}
