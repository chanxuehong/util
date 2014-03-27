// 提供一些常用的功能函数
package util

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
)

// 解析配置文件
// 配置的格式一般就是 key = value
// # 号开头的行为注释
// 函数简单的认为第一个 "=" 之前为 key, 之后为 value, 不会分析格式是否正确!
func ParseConfigFile(filename string) (map[string]string, error) {
	if filename == "" {
		return nil, errors.New("util.ParseConfigFile: filename is empty")
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		ret := make(map[string]string)
		return ret, nil
	}

	r := bufio.NewReader(bytes.NewReader(b))
	m := make(map[string]string)
	for {
		line, err := r.ReadString('\n')
		if err == nil { // 还有新行
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) < 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key == "" || value == "" {
				continue
			}
			if strings.HasPrefix(key, "#") { // 注释
				continue
			}
			m[key] = value
		} else { // 最后一行
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) < 2 {
				break
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key == "" || value == "" {
				break
			}
			if strings.HasPrefix(key, "#") { // 注释
				break
			}
			m[key] = value
			break
		}
	}
	return m, nil
}
