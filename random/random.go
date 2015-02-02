package random

import (
	"bytes"
	"crypto/sha1"
	"net"
	"os"
	"time"

	"github.com/chanxuehong/util/random/internal"
)

var (
	// 不同类型的 salt 切片 underlyingSalt 不同的部分,
	// 定期更新这个 underlyingSalt 来达到更新不同的 salt 的目的;
	// NOTE: 因为 salt 没有实际意义, 所以无需 lock.
	underlyingSalt [randomSaltLen + sessionIdSaltLen]byte

	pid            uint16   // 进程号
	realMAC        [6]byte  // 本机的某一个网卡的 MAC 地址, 如果没有则取随机数
	mac            [6]byte  // realMAC 混淆后的结果
	macSHA1HashSum [20]byte // mac 的 SHA1 哈希码
)

func init() {
	internal.ReadRandomBytes(underlyingSalt[:]) // 初始化 underlyingSalt

	go func() { // 启动一个 goroutine 定期更新 underlyingSalt
		tickChan := time.Tick(time.Minute * 5)
		for {
			select {
			case <-tickChan:
				internal.ReadRandomBytes(underlyingSalt[:])
			}
		}
	}()

	hostname, _ := os.Hostname()
	if len(hostname) < 2 {
		hostname = "hostname"
	}
	pidMask := uint16(hostname[0])<<8 | uint16(hostname[1])
	pid = uint16(os.Getpid()) ^ pidMask // 获取 pid 并混淆 pid

	realMAC = getMAC()

	// 获取 mac 并混淆, 请保证集群中所有的混淆要一致!!!
	mac = realMAC
	mac[0] ^= 0x12
	mac[1] ^= 0x34
	mac[2] ^= 0x56
	mac[3] ^= 0x78
	mac[4] ^= 0x9a
	mac[5] ^= 0xbc

	macSHA1HashSum = sha1.Sum(mac[:])
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
	internal.ReadRandomBytes(mac[:])
	mac[0] |= 0x01 // 设置多播标志, 以区分正常的 MAC
	return
}
