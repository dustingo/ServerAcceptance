/*
解析toml配置文件
*/
package util

import (
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Cpu        CpuInformation
	Interface  map[string]InterfaceInformation
	Memory     MemInformation
	Dns        DnsInformation
	System     SystemInformation
	Disk       map[string]DiskInformation
	Resolution ResolutionInformation
}

//CPU
type CpuInformation struct {
	Model     string
	Processor int
}

//system
type SystemInformation struct {
	Osname    string
	Osrelease string
	Oskernel  string
	Osarch    string
}

//网卡
type InterfaceInformation struct {
	Name  string
	Speed string
	Mtu   string
}

//dns
type DnsInformation struct {
	Nameserver []string
}

//内存
type MemInformation struct {
	Total int
}

//disk
type DiskInformation struct {
	Mountpoint string
	Size       float64
}

//解析地址需求
type ResolutionInformation struct {
	Server []string
}

//解析toml
func ParseToml(fp string) (*tomlConfig, error) {
	var tomlconfig tomlConfig
	_, err := toml.DecodeFile(fp, &tomlconfig)
	if err != nil {
		return nil, err
	}
	return &tomlconfig, nil
}
