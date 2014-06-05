package random

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sync/atomic"
	"time"
)

// 确保 Hash != nil
func commonRandom(Hash hash.Hash, localSalt []byte, otherSalts [][]byte) []byte {
	var otherSaltsTotalLen int
	for i := 0; i < len(otherSalts); i++ {
		otherSaltsTotalLen += len(otherSalts[i])
	}

	nowNanosecond := uint64(time.Now().UnixNano())

	// nowNanosecond + otherSalts + localSalt + macAddr
	allSalts := make([]byte, 8,
		8+otherSaltsTotalLen+len(localSalt)+len(macAddr))

	allSalts[0] = byte(nowNanosecond >> 56)
	allSalts[1] = byte(nowNanosecond >> 48)
	allSalts[2] = byte(nowNanosecond >> 40)
	allSalts[3] = byte(nowNanosecond >> 32)
	allSalts[4] = byte(nowNanosecond >> 24)
	allSalts[5] = byte(nowNanosecond >> 16)
	allSalts[6] = byte(nowNanosecond >> 8)
	allSalts[7] = byte(nowNanosecond)

	for i := 0; i < len(otherSalts); i++ {
		allSalts = append(allSalts, otherSalts[i]...)
	}
	allSalts = append(allSalts, localSalt...)
	allSalts = append(allSalts, macAddr...)

	Hash.Write(allSalts) // never returns an error
	return Hash.Sum(nil)
}

// The returned bytes has not been hex encoded, is raw bytes.
//  Hash = md5.New() if nil
func NewRandom(Hash hash.Hash, salts ...[]byte) []byte {
	if Hash == nil {
		Hash = md5.New()
	}
	return commonRandom(Hash, localRandomSalt, salts)
}

// The returned bytes has been hex encoded.
//  Hash = md5.New() if nil
func NewToken(Hash hash.Hash, salts ...[]byte) []byte {
	if Hash == nil {
		Hash = md5.New()
	}
	token := commonRandom(Hash, localTokenSalt, salts)
	ret := make([]byte, hex.EncodedLen(len(token)))
	hex.Encode(ret, token)
	return ret
}

// The returned bytes have been hex encoded.
func NewSessionID(salts ...[]byte) []byte {
	// 32bits unixtime + 48bits mac + 16bits pid + 32bits crc + 32bits clockSequence
	ret := make([]byte, 20)
	timenow := time.Now()

	// 写入 32bits unixtime
	nowSec := uint32(timenow.Unix())
	ret[0] = byte(nowSec >> 24)
	ret[1] = byte(nowSec >> 16)
	ret[2] = byte(nowSec >> 8)
	ret[3] = byte(nowSec)

	// 写入 48bits mac
	copy(ret[4:], macAddr)

	// 写入 16bits pid
	ret[10] = byte(pid >> 8)
	ret[11] = byte(pid)

	// 写入 clockSequence
	seq := atomic.AddUint32(&clockSequence, 1)
	ret[16] = byte(seq >> 24)
	ret[17] = byte(seq >> 16)
	ret[18] = byte(seq >> 8)
	ret[19] = byte(seq)

	// 写入 32bits crc
	var saltsTotalLen int
	for i := 0; i < len(salts); i++ {
		saltsTotalLen += len(salts[i])
	}

	nowNanosecond := uint64(timenow.UnixNano())

	// nowNanosecond + seq + salts + localSessionSalt + macAddr
	allSalts := make([]byte, 12,
		12+saltsTotalLen+len(localSessionSalt)+len(macAddr))

	allSalts[0] = byte(nowNanosecond >> 56)
	allSalts[1] = byte(nowNanosecond >> 48)
	allSalts[2] = byte(nowNanosecond >> 40)
	allSalts[3] = byte(nowNanosecond >> 32)
	allSalts[4] = byte(nowNanosecond >> 24)
	allSalts[5] = byte(nowNanosecond >> 16)
	allSalts[6] = byte(nowNanosecond >> 8)
	allSalts[7] = byte(nowNanosecond)
	allSalts[8] = byte(seq >> 24)
	allSalts[9] = byte(seq >> 16)
	allSalts[10] = byte(seq >> 8)
	allSalts[11] = byte(seq)

	for i := 0; i < len(salts); i++ {
		allSalts = append(allSalts, salts[i]...)
	}
	allSalts = append(allSalts, localSessionSalt...)
	allSalts = append(allSalts, macAddr...)

	hashSum := sha1.Sum(allSalts)
	ret[12] = hashSum[12]
	ret[13] = hashSum[13]
	ret[14] = hashSum[14]
	ret[15] = hashSum[15]

	hexRet := make([]byte, hex.EncodedLen(len(ret)))
	hex.Encode(hexRet, ret)
	return hexRet
}
