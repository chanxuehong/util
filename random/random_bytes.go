package random

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strconv"
	"time"
)

func commonRandom(newHash func() hash.Hash,
	localSalt []byte,
	otherSalts [][]byte) []byte {

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
func NewRandom(newHash func() hash.Hash, salts ...[]byte) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	return commonRandom(newHash, randomSalt, salts)
}

// The returned bytes have been hex encoded.
//
// newHash = md5.New if nil
func NewToken(newHash func() hash.Hash, salts ...[]byte) []byte {
	if newHash == nil {
		newHash = md5.New
	}
	token := commonRandom(newHash, tokenSalt, salts)
	ret := make([]byte, hex.EncodedLen(len(token)))
	hex.Encode(ret, token)
	return ret
}

// The returned bytes have been hex encoded.
//
// newHash = md5.New if nil
func NewSessionID(newHash func() hash.Hash, salts ...[]byte) []byte {
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
