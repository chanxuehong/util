package random

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/sha1"
	mathRand "math/rand"
	"net"
	"os"
	"time"
)

const (
	localSaltLen           = 32
	localSaltNum           = 3
	underlyingLocalSaltLen = localSaltLen * localSaltNum
)

var (
	// 不同类型的 localSalt 切片 underlyingLocalSalt 不同的部分,
	// 定期更新这个 underlyingLocalSalt 来达到更新不同的 localSalt 的目的;
	// NOTE: 因为 localSalt 没有实际意义, 所以无需 lock.
	underlyingLocalSalt [underlyingLocalSaltLen]byte

	// 不同的需求用不同的 localSalt, 防止暴力猜
	localRandomSalt    = underlyingLocalSalt[0*localSaltLen : 1*localSaltLen]
	localTokenSalt     = underlyingLocalSalt[1*localSaltLen : 2*localSaltLen]
	localSessionIdSalt = underlyingLocalSalt[2*localSaltLen : 3*localSaltLen]

	randomClockSequence    uint32
	idClockSequence        uint32
	uuidClockSequence      uint32
	sessionIdClockSequence uint64

	pid            uint16   // 进程号
	mac            [6]byte  // 本机的某一个网卡的 MAC 地址, 如果没有则取随机数
	macSHA1HashSum [20]byte // mac 的 SHA1 哈希码
)

// 更新 underlyingLocalSalt, 间接的更新了所有的 localSalt.
func updateUnderlyingLocalSalt() {
	// 每次获取一个字节貌似性能比较好, 虽然不是最好
	for i := 0; i < underlyingLocalSaltLen; i++ {
		n, err := cryptoRand.Prime(cryptoRand.Reader, 8)
		if err != nil {
			goto UPDATE_BY_MATH_RAND
		}
		underlyingLocalSalt[i] = n.Bytes()[0]
	}
	return //

UPDATE_BY_MATH_RAND:
	rd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < underlyingLocalSaltLen; i++ {
		underlyingLocalSalt[i] = byte(rd.Uint32())
	}
}

var zeroMAC [6]byte

// 获取一个本机的 MAC 地址, 如果没有有效的则用随机数代替.
func getMAC() (mac [6]byte) {
	interfaces, err := net.Interfaces()
	if err != nil {
		goto GEN_MAC_BY_RAND
	}

	for _, itf := range interfaces {
		if itf.Flags&net.FlagUp == net.FlagUp && // 接口是 up 的
			itf.Flags&net.FlagLoopback == 0 && // 接口不是 loopback
			len(itf.HardwareAddr) == 6 && // IEEE MAC-48, EUI-48
			!bytes.Equal(itf.HardwareAddr, zeroMAC[:]) /* 不是全0的MAC */ {

			copy(mac[:], itf.HardwareAddr)
			return
		}
	}

GEN_MAC_BY_RAND:
	copy(mac[:], localRandomSalt[localSaltLen-6:])
	mac[0] |= 0x01 // 设置多播标志, 以区分正常的 MAC
	return
}

func init() {
	updateUnderlyingLocalSalt() // 初始化 underlyingLocalSalt

	// 启动一个 goroutine 定期更新 underlyingLocalSalt
	go func() {
		tickChan := time.Tick(time.Minute * 5)
		for {
			select {
			case <-tickChan:
				updateUnderlyingLocalSalt()
			}
		}
	}()

	randomClockSequence = uint32(localRandomSalt[0])<<24 |
		uint32(localRandomSalt[1])<<16 |
		uint32(localRandomSalt[2])<<8 |
		uint32(localRandomSalt[3])

	idClockSequence = uint32(localTokenSalt[0])<<24 |
		uint32(localTokenSalt[1])<<16 |
		uint32(localTokenSalt[2])<<8 |
		uint32(localTokenSalt[3])

	uuidClockSequence = uint32(localSessionIdSalt[8])<<24 |
		uint32(localSessionIdSalt[9])<<16 |
		uint32(localSessionIdSalt[10])<<8 |
		uint32(localSessionIdSalt[11])

	sessionIdClockSequence = uint64(localSessionIdSalt[0])<<56 |
		uint64(localSessionIdSalt[1])<<48 |
		uint64(localSessionIdSalt[2])<<40 |
		uint64(localSessionIdSalt[3])<<32 |
		uint64(localSessionIdSalt[4])<<24 |
		uint64(localSessionIdSalt[5])<<16 |
		uint64(localSessionIdSalt[6])<<8 |
		uint64(localSessionIdSalt[7])

	hostname, _ := os.Hostname()
	if len(hostname) < 2 {
		hostname = "hostname"
	}
	pidMask := uint16(hostname[0])<<8 | uint16(hostname[1])
	pid = uint16(os.Getpid()) ^ pidMask // 获取 pid 并混淆 pid

	mac = getMAC()
	macSHA1HashSum = sha1.Sum(mac[:])

	// 混淆 mac;
	//  NOTE: 可以根据自己的需要来混淆, 但是集群里所有的程序 mac 都要一样的混淆
	mac[0] ^= 0x12
	mac[1] ^= 0x34
	mac[2] ^= 0x56
	mac[3] ^= 0x78
	mac[4] ^= 0x9a
	mac[5] ^= 0xbc
}
