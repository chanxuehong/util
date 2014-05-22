package util

import (
	"bytes"
	"strings"
)

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
//  NOTE: 根据 '\n' 来分行的, 如果 mac 上的 '\r' 则可以单独写一个
func TrimSpace(src []byte) []byte {
	byteSlices := bytes.Split(src, []byte{'\n'})
	for i := 0; i < len(byteSlices); i++ {
		byteSlices[i] = bytes.TrimSpace(byteSlices[i])
	}
	return bytes.Join(byteSlices, nil)
}

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
//  NOTE: 根据 '\n' 来分行的, 如果 mac 上的 '\r' 则可以单独写一个
func TrimSpaceString(src string) string {
	strs := strings.Split(src, "\n")
	for i := 0; i < len(strs); i++ {
		strs[i] = strings.TrimSpace(strs[i])
	}
	return strings.Join(strs, "")
}
