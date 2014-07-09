package security

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

// 对 content(可以为多个) 做签名(hex后的结果).
//  @salt     hash 的盐(要足够长)
//  @content  待签名的 content
//  NOTE: content 的顺序不一样结果也不一样.
func Signature(salt []byte, content ...[]byte) (hashsum []byte) {
	h := sha1.New()
	for i := 0; i < len(content); i++ {
		h.Write(content[i])
	}
	h.Write(salt)
	_hashsum := h.Sum(nil)

	hashsum = make([]byte, 32)
	hex.Encode(hashsum, _hashsum[:16])
	return
}

// 对 map[string][]byte 的 values 做签名(hex后的结果).
//  @salt    hash 的盐(要足够长)
//  @kvs     key-value pairs
func SignatureEx(salt []byte, kvs map[string][]byte) (hashsum []byte) {
	kvsLen := len(kvs)
	switch {
	case kvsLen > 1:
		// 对 kvs 的 key 做排序
		keys := make(sort.StringSlice, kvsLen)
		i := 0
		for key := range kvs {
			keys[i] = key
			i++
		}
		keys.Sort()

		values := make([][]byte, kvsLen)
		for i := 0; i < kvsLen; i++ {
			values[i] = kvs[keys[i]]
		}

		return Signature(salt, values...)

	case kvsLen == 1:
		var content []byte
		for _, value := range kvs {
			content = value
		}
		return Signature(salt, content)

	default:
		return Signature(salt)
	}
}
