package atomic

import (
	"sync/atomic"
	"testing"
)

func TestStoreLoadBool(t *testing.T) {
	var b Bool // zero value is false
	v := b.Load()
	if v {
		t.Errorf("want false, not true")
		return
	}

	b.Store(true)
	v = b.Load()
	if !v {
		t.Errorf("want true, not false")
		return
	}

	b.Store(false)
	v = b.Load()
	if v {
		t.Errorf("want false, not true")
		return
	}

	b.Store(true)
	v = b.Load()
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
		result = val.Load()
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
		val.Store(true)
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
