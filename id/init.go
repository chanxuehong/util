package id

import (
	"bytes"
	"crypto/sha1"
	"net"
	"os"

	"github.com/chanxuehong/util/random"
)

var (
	realMAC        [6]byte  // 本机的某一个网卡的 MAC 地址, 如果没有则取随机数
	fakeMAC        [6]byte  // 混淆后的 MAC
	macSHA1HashSum [20]byte // realMAC 的 SHA1 哈希码

	pid uint16 // 进程号
)

func init() {
	realMAC = getMAC()

	// 保证集群中所有的混淆要一致!!!
	fakeMAC = realMAC
	fakeMAC[0] ^= 0x12
	fakeMAC[1] ^= 0x34
	fakeMAC[2] ^= 0x56
	fakeMAC[3] ^= 0x78
	fakeMAC[4] ^= 0x9a
	fakeMAC[5] ^= 0xbc

	macSHA1HashSum = sha1.Sum(realMAC[:])

	pid = uint16(os.Getpid()) ^ 0x7788 // 保证集群中所有的混淆要一致!!!
}

var zeroMAC [6]byte

// 获取一个本机的 MAC 地址, 如果没有有效的则用随机数代替.
func getMAC() (mac [6]byte) {
	interfaces, err := net.Interfaces()
	if err != nil {
		goto GET_MAC_FROM_RAND
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

GET_MAC_FROM_RAND:
	random.ReadRandomBytes(mac[:])
	mac[0] |= 0x01 // 设置多播标志, 以区分正常的 MAC
	return
}
