//预检
package work

import (
	"fmt"
	"log"
	"math"

	"github.com/dustingo/ServerAcceptance/util"
)

// PreCheck 预检查函数
func PreCheck(url string) {
	//下载配置文件
	//confURL := "https://raw.githubusercontent.com/dustingo/configtoml/main/config.toml"
	wgetErr := Wget(url, "config.toml")
	if wgetErr != nil {
		fmt.Println(wgetErr)
		return
	}

	//cpu
	c := new(util.CpuInfo)
	c.GetCpu()
	cModuleName := c.ModelName
	cProcessor := c.Processor
	//system
	s := new(util.SystemInfo)
	s.GetSystem()
	sName := s.OSName
	sRelease := s.OSRelease
	sKernel := s.OSKernel
	sArch := s.OSArch
	//memory
	mTotal, _, _, memerr := util.GetMemory()
	if memerr != nil {
		log.Fatalln("Get memory error:", memerr)
	}
	//disk
	diskStats, diskerr := util.GetDiskStats()
	if diskerr != nil {
		fmt.Println(diskerr)
	}
	//dns
	dns, dnserr := util.GetDNSInfo()
	if dnserr != nil {
		log.Fatalln("Get dns Error:", dnserr)
	}

	//

	//解析配置文件
	tc, tcerr := util.ParseToml("config.toml")
	if tcerr != nil {
		log.Fatalln(tcerr)
		return
	}
	//检查system
	fmt.Println("SYSTEM:")
	if tc.System.Osname == sName {
		fmt.Printf("OS Name: %s [ok]\n", tc.System.Osname)
	} else {
		fmt.Printf("OS Name: %s [ERR]\n", sName)
		fmt.Printf("OS Name Config: %s \n", tc.System.Osname)
	}

	if tc.System.Osrelease == sRelease {
		fmt.Printf("OS Release: %s [OK]\n", tc.System.Osrelease)
	} else {
		fmt.Printf("OS Release: %s [ERR]\n", sRelease)
		fmt.Printf("OS Release Config: %s \n", tc.System.Osrelease)
	}
	if tc.System.Oskernel == sKernel {
		fmt.Printf("OS Kernel: %s [OK]\n", tc.System.Oskernel)
	} else {
		fmt.Printf("OS Kernel: %s [ERR]\n", sKernel)
		fmt.Printf("OS Kernel Config: %s \n", tc.System.Oskernel)
	}
	if tc.System.Osarch == sArch {
		fmt.Printf("OS Arch: %s [OK]\n", tc.System.Osarch)
	} else {
		fmt.Printf("OS Arch: %s [ERR]\n", sArch)
		fmt.Printf("OS Arch Config: %s \n", tc.System.Osarch)
	}
	//检查CPU
	fmt.Println("CPU:")
	if tc.Cpu.Model == cModuleName {
		fmt.Printf("Module Name: %s [OK]\n", tc.Cpu.Model)
	} else {
		fmt.Printf("Module Name: %s [ERR]\n", cModuleName)
		fmt.Printf("Module Name config: %s\n", tc.Cpu.Model)
	}
	if tc.Cpu.Processor == cProcessor {
		fmt.Printf("Processor Num: %d [OK]\n", tc.Cpu.Processor)
	} else {
		fmt.Printf("Processor Num: %d [ERR]\n", cProcessor)
		fmt.Printf("Processor Num Config: %d\n", tc.Cpu.Processor)
	}
	//检查内存,误差控制在1024KB
	fmt.Println("MEMORY:")
	if math.Abs(float64(tc.Memory.Total-mTotal)) < 1024 {
		fmt.Printf("Mem Total: %d [ok]\n", tc.Memory.Total)
	} else {
		fmt.Printf("Mem Total: %d [ERR]\n", mTotal)
		fmt.Printf("Mem Total Config: %d \n", tc.Memory.Total)
	}
	//检查硬盘
	fmt.Println("Disk:")
	diskMap := make(map[string]float64)
	for _, stat := range diskStats {
		diskMap[stat.Labels.MountPoint] = stat.Size
	}
	for mp := range tc.Disk {
		if _, ok := diskMap[tc.Disk[mp].Mountpoint]; ok {
			//磁盘误差在100MB
			if diskMap[tc.Disk[mp].Mountpoint]-tc.Disk[mp].Size < 100.0 {
				fmt.Printf("Disk MountPoint  %s Size: %f [ok]\n", tc.Disk[mp].Mountpoint, diskMap[tc.Disk[mp].Mountpoint])
			} else {
				fmt.Printf("Disk MountPoint  %s Size: %f [ERR]\n", tc.Disk[mp].Mountpoint, diskMap[tc.Disk[mp].Mountpoint])
				fmt.Printf("Disk MountPoint   %s Config Size: %f \n", tc.Disk[mp].Mountpoint, tc.Disk[mp].Size)
			}
		} else {
			fmt.Printf("Disk MountPoint %s not exists\n", tc.Disk[mp].Mountpoint)
		}
	}
	//检查网卡
	fmt.Println("InterFace:")
	for k := range tc.Interface {
		speed, ip, mtu, neterr := util.GetNetInfo(k)
		if neterr != nil {
			fmt.Printf("Get ifcfg-%s ERROR\n", k)
			continue
		}
		if speed == tc.Interface[k].Speed {
			fmt.Printf("NET ifcfg-%s Speed: %s [ok]\n", k, speed)
		} else {
			fmt.Printf("NET ifcfg-%s Speed: %s [ERR]\n", k, speed)
		}
		if ip == "" {
			fmt.Printf("NET ifcfg-%s ip is null [ERR]\n", k)
		} else {
			fmt.Printf("NET ifcfg-%s ip is: %s \n", k, ip)
		}
		if mtu == tc.Interface[k].Mtu {
			fmt.Printf("NET ifcfg-%s MTU:%s [OK]\n", k, mtu)
		} else {
			fmt.Printf("NET ifcfg-%s MTU:%s [ERR]\n", k, mtu)
		}

	}
	//检查dns
	fmt.Println("DNS:")
	dnsServerMap := make(map[string]int)
	lostDNS := []string{}
	for n, d := range dns {
		dnsServerMap[d] = n

	}
	for _, configDNS := range tc.Dns.Nameserver {
		if _, ok := dnsServerMap[configDNS]; ok {
			continue
		}
		lostDNS = append(lostDNS, configDNS)
	}

	if len(lostDNS) == 0 {
		fmt.Printf("DNS Info: %v [OK]\n", tc.Dns.Nameserver)
	} else {
		fmt.Printf("DNS Info: %v [ERR]\n", dns)
		fmt.Printf("DNS LOST: %v \n", lostDNS)
	}
}
