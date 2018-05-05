package atomic

import (
	"sync/atomic"
	"testing"
)

func TestStoreLoadBool(t *testing.T) {
	var b Bool // zero value is false
	v := LoadBool(&b)
	if v {
		t.Errorf("want false, not true")
		return
	}

	StoreBool(&b, true)
	v = LoadBool(&b)
	if !v {
		t.Errorf("want true, not false")
		return
	}

	StoreBool(&b, false)
	v = LoadBool(&b)
	if v {
		t.Errorf("want false, not true")
		return
	}

	StoreBool(&b, true)
	v = LoadBool(&b)
	if !v {
		t.Errorf("want true, not false")
		return
	}
}

func BenchmarkLoadBool(b *testing.B) {
	var val Bool
	var result bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = LoadBool(&val)
	}
	_ = result
}

func BenchmarkLoadAtomicValue(b *testing.B) {
	var val atomic.Value
	val.Store(false)
	var result bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = val.Load().(bool)
	}
	_ = result
}

func BenchmarkStoreBool(b *testing.B) {
	var val Bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StoreBool(&val, true)
	}
}

func BenchmarkStoreAtomicValue(b *testing.B) {
	var val atomic.Value
	val.Store(false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Store(true)
	}
}
