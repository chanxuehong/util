package security

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"sort"
)

// 对 content 做签名(hex后的结果).
// @salt     hash 的盐
// @content  待签名的 content
func Signature(salt []byte, content ...[]byte) []byte {
	h := hmac.New(md5.New, salt)
	for i := 0; i < len(content); i++ {
		h.Write(content[i])
	}
	hashSum := h.Sum(nil)

	ret := make([]byte, hex.EncodedLen(len(hashSum)))
	hex.Encode(ret, hashSum)
	return ret
}

// 对 map[string][]byte 的 values 做签名(hex后的结果).
// @salt    hash 的盐
// @kvs     key-value pairs
func MapSignature(salt []byte, kvs map[string][]byte) []byte {
	switch {
	case len(kvs) > 1:
		// 对 kvs 的 key 做排序
		keys := make(sort.StringSlice, len(kvs))
		i := 0
		for key := range kvs {
			keys[i] = key
			i++
		}
		keys.Sort()

		values := make([][]byte, len(keys))
		for i := 0; i < len(keys); i++ {
			values[i] = kvs[keys[i]]
		}

		return Signature(salt, values...)

	case len(kvs) == 1:
		var content []byte
		for _, value := range kvs {
			content = value
		}
		return Signature(salt, content)

	default:
		return Signature(salt)
	}
}
