package random

import (
	cryptoRand "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"time"
)

const (
	saltLen       = 128
	saltNum       = 3
	bottomSaltLen = saltLen * saltNum
)

var (
	// 实际上只用定期更新这个 bottomSalt 即可.
	bottomSalt [bottomSaltLen]byte
	// 不同的需求使用不同的 salt, 防止暴力猜
	randomSalt  = bottomSalt[0*saltLen : 1*saltLen]
	tokenSalt   = bottomSalt[1*saltLen : 2*saltLen]
	sessionSalt = bottomSalt[2*saltLen : 3*saltLen]
)

func updateSalt() {
	// 优先用 crypto/rand 生成 salts
	var bi *big.Int
	var err error
	for i := 0; i < bottomSaltLen; i++ {
		// 每次获取一个字节貌似性能比较好, 虽然不是最好
		bi, err = cryptoRand.Prime(cryptoRand.Reader, 8)
		if err != nil {
			goto MATH_RAND
		}
		bottomSalt[i] = bi.Bytes()[0]
	}

MATH_RAND: // crypto/rand 生成 failed, 就用 math/rand 来生成 salts
	rd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < bottomSaltLen; i++ {
		bottomSalt[i] = byte(rd.Uint32())
	}
}

func init() {
	// 初始化 salt
	updateSalt()
	// 每5分钟更新一次 salt
	go func() {
		ch := time.Tick(time.Minute * 5)
		for {
			select {
			case <-ch:
				updateSalt()
			}
		}
	}()
}
