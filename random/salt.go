package random

import (
	cryptoRand "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"os"
	"strconv"
	"time"
)

const (
	localSaltLen           = 128
	localSaltNum           = 3
	underlyingLocalSaltLen = localSaltLen * localSaltNum
)

var (
	pid      string
	hostname string

	// 不同的需求用不同的 local salt, 防止暴力猜. 所有的这些 local salt 切片一个底层的数组
	// underlyingLocalSalt 的不同部分, 定期更新这个 underlyingLocalSalt 来达到更新不同的 local salt;
	// NOTE: 因为 local salt 没有实际意义, 所以无需 lock.
	underlyingLocalSalt [underlyingLocalSaltLen]byte
	localRandomSalt     = underlyingLocalSalt[0*localSaltLen : 1*localSaltLen]
	localTokenSalt      = underlyingLocalSalt[1*localSaltLen : 2*localSaltLen]
	localSessionSalt    = underlyingLocalSalt[2*localSaltLen : 3*localSaltLen]
)

func updateUnderlyingLocalSalt() {
	// 优先用 crypto/rand
	var bi *big.Int
	var err error
	for i := 0; i < underlyingLocalSaltLen; i++ {
		bi, err = cryptoRand.Prime(cryptoRand.Reader, 8) // 每次获取一个字节貌似性能比较好, 虽然不能保证最好
		if err != nil {
			goto MATH_RAND
		}
		underlyingLocalSalt[i] = bi.Bytes()[0]
	}

	return // crypto/rand 更新成功, 返回

	// crypto/rand 更新失败, 就用 math/rand 来更新
MATH_RAND:
	rd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < underlyingLocalSaltLen; i++ {
		underlyingLocalSalt[i] = byte(rd.Uint32())
	}
}

func init() {
	pid = strconv.FormatUint(uint64(os.Getpid()), 16)
	hostname, _ = os.Hostname()

	updateUnderlyingLocalSalt() // 初始化 salt

	go func() {
		ch := time.Tick(time.Minute * 5) // 每5分钟更新一次 salt
		for {
			select {
			case <-ch:
				updateUnderlyingLocalSalt()
			}
		}
	}()
}
