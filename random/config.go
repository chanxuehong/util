package random

import (
	"bytes"
	cryptoRand "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"net"
	"os"
	"time"
)

const (
	localSaltLen           = 128
	localSaltNum           = 3
	underlyingLocalSaltLen = localSaltLen * localSaltNum
)

var (
	macAddr []byte // 本机的一个网卡的 MAC 地址, 如果没有则取随机数
	pid     uint16 // 进程号

	// 类似 uuid 里的 clockSequence
	randomClockSequence  uint32
	tokenClockSequence   uint32
	sessionClockSequence uint32

	// 不同的需求用不同的 local salt, 防止暴力猜. 所有的这些 local salt 切片一个底层的数组
	// underlyingLocalSalt 的不同部分, 定期更新这个 underlyingLocalSalt 来达到更新不同的 local salt;
	// NOTE: 因为 local salt 没有实际意义, 所以无需 lock.
	underlyingLocalSalt [underlyingLocalSaltLen]byte
	localRandomSalt     = underlyingLocalSalt[0*localSaltLen : 1*localSaltLen]
	localTokenSalt      = underlyingLocalSalt[1*localSaltLen : 2*localSaltLen]
	localSessionSalt    = underlyingLocalSalt[2*localSaltLen : 3*localSaltLen]
)

// 更新底层 salt 数组, 间接的更新了所有的 salt
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

// 获取一个本机的 MAC 地址, 如果没有有效的则用随机数代替
func getHardwareAddress() []byte {
	itfs, err := net.Interfaces()
	if err != nil {
		goto generate_mac
	}
	for _, itf := range itfs {
		if itf.Flags&net.FlagUp != 0 && // 接口是 up 的
			itf.Flags&net.FlagLoopback == 0 && // 接口不是 loopback
			len(itf.HardwareAddr) == 6 && // IEEE MAC-48, EUI-48
			!bytes.Equal(itf.HardwareAddr, make([]byte, len(itf.HardwareAddr))) { // 不是全 0 地址

			return itf.HardwareAddr
		}
	}

generate_mac:
	// 没有找到有效的 MAC 地址, 只能随机生成 MAC 地址了;
	// 这里直接用 localRandomSalt 的前 6 位了
	mac := localRandomSalt[:6]
	mac[0] |= 0x01 // 设置多播标志, 以区分正常的 MAC
	return mac
}

func init() {
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

	macAddr = getHardwareAddress()
	pid = uint16(os.Getpid())

	mathRand.Seed(time.Now().UnixNano())
	randomClockSequence = mathRand.Uint32()
	tokenClockSequence = mathRand.Uint32()
	sessionClockSequence = mathRand.Uint32()
}
