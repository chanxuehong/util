package version

import (
	"fmt"
	"strconv"
	"strings"
)

var _ fmt.Stringer = Version{}

// Version 表示一个版本号.
type Version struct {
	Major, Minor, Patch int
}

func NewVersion(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// String 将 Version 格式化成 x.y.z 的字符串形式.
func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// Compare 比较 v, v2 的大小.
//  返回 -1 表示 v < v2
//  返回 0  表示 v == v2
//  返回 +1  表示 v > v2
func (v Version) Compare(v2 Version) int {
	return Compare(v, v2)
}

// Parse 解析字符串 x.y.z 到 Version 对象.
func Parse(str string) (v Version, err error) {
	var (
		strCopy  = str // 用于错误显示
		dotIndex int
	)

	// 获取 Major
	dotIndex = strings.IndexByte(str, '.')
	switch {
	case dotIndex > 0:
		v.Major, err = strconv.Atoi(str[:dotIndex])
		if err != nil {
			err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
			return
		}
		str = str[dotIndex+1:]
	case dotIndex == 0:
		err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
		return
	case dotIndex < 0:
		v.Major, err = strconv.Atoi(str)
		if err != nil {
			err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
			return
		}
		return // 没有更多的 '.' 了, 直接返回
	}

	// 获取 Minor
	dotIndex = strings.IndexByte(str, '.')
	switch {
	case dotIndex > 0:
		v.Minor, err = strconv.Atoi(str[:dotIndex])
		if err != nil {
			err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
			return
		}
		str = str[dotIndex+1:]
	case dotIndex == 0:
		err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
		return
	case dotIndex < 0:
		v.Minor, err = strconv.Atoi(str)
		if err != nil {
			err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
			return
		}
		return // 没有更多的 '.' 了, 直接返回
	}

	// 获取 Patch
	dotIndex = strings.IndexByte(str, '.')
	if dotIndex >= 0 {
		err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
		return
	}
	v.Patch, err = strconv.Atoi(str)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal string %q into Go value of type Version", strCopy)
		return
	}
	return
}
