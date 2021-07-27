/*
网卡信息，返回网卡速率，ip地址，mtu，error
*/
package util

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

type NFError struct {
	msg string
}

func (nf *NFError) Error() string {
	return nf.msg
}

// 获取网卡信息
func GetNetInfo(name string) (sp string, ipa string, mtu string, error *NFError) {
	var ip string
	var out bytes.Buffer
	var speed string
	netFace, err := net.InterfaceByName(name)
	if err != nil {
		return "", "", "", &NFError{err.Error()}
	}
	//如果当前网卡处于未up状态，直接返回错误
	//if !strings.Contains(netFace.Flags.String(), "up") {

	//	return "", "", &NFError{fmt.Sprintf("interface %s did not up", name)}
	//}
	ips, _ := netFace.Addrs()
	if len(ips) == 1 {
		ip = ips[0].String()
	} else {

		return "", "", "", &NFError{fmt.Sprintf("interface %s has no ipadress", name)}
	}
	cmd := exec.Command("ethtool", name)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", "", "", &NFError{err.Error()}
	}
	ethtoolOutSlice := strings.Split(strings.Replace(out.String(), " ", "", -1), "\n")
	for _, ethtoolOut := range ethtoolOutSlice {
		if strings.HasPrefix(strings.TrimSpace(ethtoolOut), "Speed") {
			speed = strings.Split(ethtoolOut, ":")[1]
		}
	}
	if speed == "" {

		return "", "", "", &NFError{fmt.Sprintf("%s speed not found\n", name)}
	}

	return speed, ip, strconv.Itoa(netFace.MTU), nil
}
