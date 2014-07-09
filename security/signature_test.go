package security

import (
	"bytes"
	"testing"
)

func TestSignature(t *testing.T) {
	Signature(nil) // 检查极端环境下是否 panic

	var salt = []byte("salt")

	content0 := []byte("content0")
	content1 := []byte("content1")
	content2 := []byte("content2")
	content3 := []byte("content0content1content2")

	hashSum0 := Signature(salt, content0, content1, content2)
	hashSum1 := Signature(salt, content3)

	if !bytes.Equal(hashSum0, hashSum1) {
		t.Error("test Signature() failed")
	}
}

func TestSignatureEx(t *testing.T) {
	m := map[string][]byte{
		"b": []byte("valueB"),
		"a": []byte("valueA"),
		"c": []byte("valueC"),
	}

	salt := []byte("salt")
	hashSum0 := Signature(salt, []byte("valueAvalueBvalueC"))
	hashSum1 := SignatureEx(salt, m)

	if !bytes.Equal(hashSum0, hashSum1) {
		t.Error("test SignatureEx() failed")
	}
}
