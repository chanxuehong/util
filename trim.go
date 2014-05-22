package util

import (
	"bytes"
	"strings"
)

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
//  NOTE: 根据 '\n' 来分行的, 如果 mac 上的 '\r' 则可以单独写一个
func TrimSpace(src []byte) []byte {
	byteSlices := bytes.Split(src, []byte{'\n'})
	for i, byteSlicesLen := 0, len(byteSlices); i < byteSlicesLen; i++ {
		byteSlices[i] = bytes.TrimSpace(byteSlices[i])
	}
	return bytes.Join(byteSlices, nil)
}

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
//  NOTE: 根据 '\n' 来分行的, 如果 mac 上的 '\r' 则可以单独写一个
func TrimSpaceString(src string) string {
	strs := strings.Split(src, "\n")
	for i, strsLen := 0, len(strs); i < strsLen; i++ {
		strs[i] = strings.TrimSpace(strs[i])
	}
	return strings.Join(strs, "")
}
