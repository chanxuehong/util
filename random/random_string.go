package random

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strconv"
	"time"
)

func commonRandomString(newHash func() hash.Hash, localSalt []byte, otherSalts []string) []byte {
	var otherSaltsTotalLen int
	for i := 0; i < len(otherSalts); i++ {
		otherSaltsTotalLen += len(otherSalts[i])
	}

	nowNanosecond := strconv.FormatUint(uint64(time.Now().UnixNano()), 16)

	// nowNanosecond + otherSalts + localSalt + pid + hostname
	allSalts := make([]byte, len(nowNanosecond),
		len(nowNanosecond)+otherSaltsTotalLen+len(localSalt)+len(pid)+len(hostname))

	copy(allSalts, nowNanosecond)
	for i := 0; i < len(otherSalts); i++ {
		allSalts = append(allSalts, otherSalts[i]...)
	}
	allSalts = append(allSalts, localSalt...)
	allSalts = append(allSalts, pid...)
	allSalts = append(allSalts, hostname...)

	h := newHash()
	h.Write(allSalts) // never returns an error.

	return h.Sum(nil)
}

// The returned bytes has not been hex encoded, is raw bytes.
//  newHash = md5.New if nil
func NewRandomString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	return commonRandomString(newHash, localRandomSalt, salts)
}

// The returned bytes has been hex encoded.
//  newHash = md5.New if nil
func NewTokenString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	token := commonRandomString(newHash, localTokenSalt, salts)
	ret := make([]byte, hex.EncodedLen(len(token)))
	hex.Encode(ret, token)
	return ret
}

// The returned bytes have been hex encoded.
//  newHash = md5.New if nil
func NewSessionIDString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}

	var saltsTotalLen int
	for i := 0; i < len(salts); i++ {
		saltsTotalLen += len(salts[i])
	}

	nowNanosecond := strconv.FormatUint(uint64(time.Now().UnixNano()), 16)

	// nowNanosecond + salts + localSessionSalt + pid + hostname
	allSalts := make([]byte, len(nowNanosecond),
		len(nowNanosecond)+saltsTotalLen+len(localSessionSalt)+len(pid)+len(hostname))

	copy(allSalts, nowNanosecond)
	for i := 0; i < len(salts); i++ {
		allSalts = append(allSalts, salts[i]...)
	}
	allSalts = append(allSalts, localSessionSalt...)
	allSalts = append(allSalts, pid...)
	allSalts = append(allSalts, hostname...)

	h := newHash()
	h.Write(allSalts) // never returns an error.

	hashsum := h.Sum(nil)
	ret := make([]byte, len(nowNanosecond)+hex.EncodedLen(len(hashsum)))
	copy(ret, nowNanosecond)
	hex.Encode(ret[len(nowNanosecond):], hashsum)
	return ret
}
