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
	localSaltLen           = 32 // 每个 salt 的长度为, 单位为字节
	localSaltNum           = 3
	underlyingLocalSaltLen = localSaltLen * localSaltNum
)

var (
	pid                    uint16 // 进程号
	macAddr                []byte // 本机的某一个网卡的 MAC 地址, 如果没有则取随机数
	macAddrHashSum         []byte // macAddr 的 SHA-1 结果的前 3 个字节
	sessionIdClockSequence uint64 // 类似 uuid 里的 clockSequence
	idClockSequence        uint32 // 类似 uuid 里的 clockSequence

	// 不同的需求用不同的 local salt, 防止暴力猜. 所有的这些 local salt 切片一个底层的数组
	// underlyingLocalSalt 的不同部分
	localRandomSalt  = underlyingLocalSalt[0*localSaltLen : 1*localSaltLen]
	localTokenSalt   = underlyingLocalSalt[1*localSaltLen : 2*localSaltLen]
	localSessionSalt = underlyingLocalSalt[2*localSaltLen : 3*localSaltLen]

	// 定期更新这个 underlyingLocalSalt 来达到更新不同的 local salt;
	// NOTE: 因为 local salt 没有实际意义, 所以无需 lock.
	underlyingLocalSalt [underlyingLocalSaltLen]byte
)

// 更新底层 salt 数组, 间接的更新了所有的 salt
func updateUnderlyingLocalSalt() {
	// 优先用 crypto/rand
	// 每次获取一个字节貌似性能比较好, 虽然不是最好
	for i := 0; i < underlyingLocalSaltLen; i++ {
		n, err := cryptoRand.Prime(cryptoRand.Reader, 8)
		if err != nil {
			goto MATH_RAND_UPDATE
		}
		underlyingLocalSalt[i] = n.Bytes()[0]
	}

	return //

	// crypto/rand 更新失败, 就用 math/rand 来更新
MATH_RAND_UPDATE:
	rd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < underlyingLocalSaltLen; i++ {
		underlyingLocalSalt[i] = byte(rd.Uint32())
	}
}

// 获取一个本机的 MAC 地址, 如果没有有效的则用随机数代替
func getHardwareAddress() []byte {
	interfaces, err := net.Interfaces()
	if err != nil {
		goto GEN_MAC_BY_RAND
	}
	for _, itf := range interfaces {
		if itf.Flags&net.FlagUp == net.FlagUp && // 接口是 up 的
			itf.Flags&net.FlagLoopback == 0 && // 接口不是 loopback
			len(itf.HardwareAddr) == 6 && // IEEE MAC-48, EUI-48
			!bytes.Equal(itf.HardwareAddr, make([]byte, 6)) /* 不是全0的MAC */ {

			return itf.HardwareAddr
		}
	}

GEN_MAC_BY_RAND: // 没有找到有效的 MAC 地址, 只能随机生成 MAC 地址了

	// 这里直接用 localRandomSalt 的后 6 位了
	mac := make([]byte, 6)
	copy(mac, localRandomSalt[localSaltLen-6:]) // localRandomSalt 会变化, 所以要 copy
	mac[0] |= 0x01                              // 设置多播标志, 以区分正常的 MAC
	return mac
}

func init() {
	updateUnderlyingLocalSalt() // 初始化 salt
	go func() {
		tickChan := time.Tick(time.Minute * 5) // 每5分钟更新一次 salt
		for {
			select {
			case <-tickChan:
				updateUnderlyingLocalSalt()
			}
		}
	}()

	hostname, _ := os.Hostname()
	if len(hostname) < 2 {
		hostname = "hostname"
	}
	hostnameBytes := []byte(hostname)
	pidMask := uint16(hostnameBytes[0])<<8 | uint16(hostnameBytes[1])
	pid = uint16(os.Getpid()) ^ pidMask // 获取 pid 并混淆 pid

	macAddr = getHardwareAddress()
	// 混淆 macAddr;
	//  NOTE: 可以根据自己的需要来混淆, 但是集群里所有的程序 macAddr 都要一样的混淆
	macAddr[0] ^= 0x12
	macAddr[1] ^= 0x34
	macAddr[2] ^= 0x56
	macAddr[3] ^= 0x78
	macAddr[4] ^= 0x9a
	macAddr[5] ^= 0xbc

	macAddrHashSumArray := sha1.Sum(macAddr)
	macAddrHashSum = macAddrHashSumArray[:3]

	sessionIdClockSequence = uint64(localSessionSalt[0])<<56 +
		uint64(localSessionSalt[1])<<48 +
		uint64(localSessionSalt[2])<<40 +
		uint64(localSessionSalt[3])<<32 +
		uint64(localSessionSalt[4])<<24 +
		uint64(localSessionSalt[5])<<16 +
		uint64(localSessionSalt[6])<<8 +
		uint64(localSessionSalt[7])

	idClockSequence = uint32(localSessionSalt[8])<<24 +
		uint32(localSessionSalt[9])<<16 +
		uint32(localSessionSalt[10])<<8 +
		uint32(localSessionSalt[11])
}
