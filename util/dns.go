/*
dns信息,返回dns地址字符串数组，错误
*/
package util

import (
	"bufio"
	"os"
	"strings"
)

// 自定义dns错误
type dnsError struct {
	msg string
}

func (d *dnsError) Error() string {
	return d.msg
}

// 获取dns信息
func GetDNSInfo() ([]string, error) {
	var nameserver []string
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "nameserver") {
			nameserver = append(nameserver, strings.Fields(line)[1])
		}
	}
	if len(nameserver) == 0 {
		return nil, &dnsError{"not found dns server"}
	}
	return nameserver, nil
}
