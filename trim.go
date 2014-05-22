package util

import (
	"regexp"
)

var trimSpaceRegexp *regexp.Regexp

func init() {
	// 匹配文本开头和结尾空白, 或者换行符极其换行符两端的空白
	trimSpaceRegexp = regexp.MustCompile(`(^\s+)|(\s+$)|([\x20\t\f]*(((\r\n)|\n|\r)[\x20\t\f]*)+)`)
}

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
func TrimSpace(src []byte) []byte {
	return trimSpaceRegexp.ReplaceAllLiteral(src, nil)
}

// 去掉 src 开头和结尾的空白, 如果 src 包括换行, 去掉换行和这个换行符两边的空白
func TrimSpaceString(src string) string {
	return trimSpaceRegexp.ReplaceAllLiteralString(src, "")
}
