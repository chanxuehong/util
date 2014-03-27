package random

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strconv"
	"time"
)

func commonRandomString(newHash func() hash.Hash,
	localSalt []byte,
	otherSalts []string) []byte {

	var otherSaltsTotalLen int
	for i := 0; i < len(otherSalts); i++ {
		otherSaltsTotalLen += len(otherSalts[i])
	}

	nowNanosecond := strconv.FormatInt(time.Now().UnixNano(), 16)
	// unix nanosecond + otherSalts + localSalt
	allSalts := make([]byte, len(nowNanosecond),
		len(nowNanosecond)+otherSaltsTotalLen+len(localSalt))
	copy(allSalts, nowNanosecond)
	for i := 0; i < len(otherSalts); i++ {
		allSalts = append(allSalts, otherSalts[i]...)
	}
	allSalts = append(allSalts, localSalt...)

	h := newHash()
	h.Write(allSalts) // never returns an error.

	return h.Sum(nil)
}

// newHash = md5.New if nil
func NewRandomString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	return commonRandomString(newHash, randomSalt, salts)
}

// The returned bytes have been hex encoded.
//
// newHash = md5.New if nil
func NewTokenString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	token := commonRandomString(newHash, tokenSalt, salts)
	ret := make([]byte, hex.EncodedLen(len(token)))
	hex.Encode(ret, token)
	return ret
}

// The returned bytes have been hex encoded.
//
// newHash = md5.New if nil
func NewSessionIDString(newHash func() hash.Hash, salts ...string) []byte {
	if newHash == nil {
		newHash = md5.New
	}

	var saltsTotalLen int
	for i := 0; i < len(salts); i++ {
		saltsTotalLen += len(salts[i])
	}

	timenow := time.Now()
	nowSecend := strconv.FormatUint(uint64(timenow.Unix()), 16)
	nowNanosecond := strconv.FormatInt(timenow.UnixNano(), 16)
	// unix nanosecond + salts + sessionSalt
	allSalts := make([]byte, len(nowNanosecond),
		len(nowNanosecond)+saltsTotalLen+len(sessionSalt))
	copy(allSalts, nowNanosecond)
	for i := 0; i < len(salts); i++ {
		allSalts = append(allSalts, salts[i]...)
	}
	allSalts = append(allSalts, sessionSalt...)

	h := newHash()
	h.Write(allSalts) // never returns an error.

	hashsum := h.Sum(nil)
	ret := make([]byte, len(nowSecend)+hex.EncodedLen(len(hashsum)))
	copy(ret, nowSecend)
	hex.Encode(ret[len(nowSecend):], hashsum)
	return ret
}
