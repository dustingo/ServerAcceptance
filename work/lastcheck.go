package work

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/dustingo/ServerAcceptance/util"
	"github.com/lorenzosaino/go-sysctl"
)

var wg sync.WaitGroup
var lconfig util.LastCheckInfo

func LastCheck(url string) {
	err := Wget(url, "lastcheck.toml")
	if err != nil {
		fmt.Println(err)
		return
	}
	last, err := util.LastParse("lastcheck.toml")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Service
	offState, offLabel, offName, onState, onLabel, onName := last.GetService()
	// Package
	yumState, yumLabe, yumName, pipState, pipLabel, pipName, perlState, perlLabel, perlName := last.GetPackage()
	// Directory
	dirInfo := last.GetDirectory()
	// Ulimit
	ulimit := last.GetUlimit()
	// syskernel
	syskernel := last.GetSysKernel()
	// dns
	dnsinfo := last.GetDNS()
	wg.Add(9)
	go packageCMD(yumState, yumName, yumLabe)
	go packageCMD(pipState, pipName, pipLabel)
	go packageCMD(perlState, perlName, perlLabel)
	go serviceCMD(offState, offName, offLabel)
	go serviceCMD(onState, onName, onLabel)
	go directoryCMD(dirInfo)
	go ulimitCMD(ulimit)
	go sysKernelCMD(syskernel)
	go dnsCMD(dnsinfo)
	wg.Wait()
}

//  serviceCMD 校验Service模块
func serviceCMD(status int, names []string, label string) {
	defer wg.Done()
	if status == 1 {
		cmd := exec.Command("systemctl", "list-units")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		if label == "off" {
			scanner := bufio.NewScanner(&out)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) == 0 {
					continue
				}
				for _, unit := range names {
					if strings.Contains(line, unit) {
						if strings.Fields(line)[3] == "active" || strings.Fields(line)[3] == "running" {
							fmt.Printf("%s still running \n", unit)
						}
					}

				}
			}

		} else if label == "on" {
			for _, unit := range names {
				if strings.Contains(out.String(), unit) {
					scanner := bufio.NewScanner(&out)
					for scanner.Scan() {
						line := scanner.Text()
						if len(line) == 0 {
							continue
						}
						if strings.Contains(line, unit) {
							if strings.Fields(line)[3] == "active" || strings.Fields(line)[3] == "running" {
								continue
							} else {
								fmt.Printf("%s not running \n", unit)
							}
						}
					}
				} else {
					fmt.Printf("%s not running \n", unit)
				}
			}
		}
	} else {
		fmt.Printf("Service.%s state is %d ignored\n", label, status)
	}
}

// packageCMD 校验Package 模块
func packageCMD(status int, names []string, label string) {
	defer wg.Done()
	if status == 1 {
		switch label {
		case "yum":
			for _, pack := range names {
				var out bytes.Buffer
				cmd := exec.Command("rpm", "-q", pack)
				cmd.Stdout = &out
				cmd.Run()
				fmt.Printf(out.String())
			}
		case "pip3", "pip":
			for _, pack := range names {
				var out bytes.Buffer
				cmd := exec.Command(label, "show", pack)
				cmd.Stdout = &out
				cmd.Run()
				if !strings.Contains(out.String(), "Name") {
					fmt.Printf("package %s is not installed \n", pack)
				}
			}
		case "perl":
			for _, pack := range names {
				var out bytes.Buffer
				cmd := exec.Command(label, "-le", "use %s", pack)
				cmd.Stderr = &out
				cmd.Run()
				if len(out.String()) != 0 {
					fmt.Printf("package %s is not installed \n", pack)
				}
			}
		default:
			fmt.Println("unknown type of method ", label)
		}
	} else {
		fmt.Printf("Package.%s status is %d ignored\n", label, status)
	}
}

// directoryCMD 校验目录模块
func directoryCMD(dirMap map[string]util.DirInfo) {
	defer wg.Done()
	for k, _ := range dirMap {
		if dirMap[k].State != 1 {
			fmt.Printf("Directory or File %s state is %d ignored \n", dirMap[k].Path, dirMap[k].State)
			continue
		}
		f, err := dirModeCheck(dirMap[k].Path)
		if err != nil {
			fmt.Println(errors.New(fmt.Sprintf("Directory or File %s not exists", dirMap[k].Path)))
			return
		}
		modEight, _ := strconv.ParseInt(fmt.Sprintf("%04o", f.Mode().Perm()), 10, 64)
		if modEight != dirMap[k].Mode {
			fmt.Printf("Directory or File %s perm is %d \n", dirMap[k].Path, modEight)
			return
		}
		uid := f.Sys().(*syscall.Stat_t).Uid
		currUser, _ := user.LookupId(fmt.Sprint(uid))
		if currUser.Username != dirMap[k].Owner {
			fmt.Printf("Directory or File %s owner is %s \n", dirMap[k].Path, currUser.Username)
			return
		}
	}
}

// ulimitCMD 比较ulimit信息，以/etc/security/limits.conf 配置文件为准
func ulimitCMD(u *util.Ulimits) {
	defer wg.Done()
	var confLimit []string
	coreSoft := strings.Join([]string{u.Core.Soft.Domain, u.Core.Soft.Type, u.Core.Soft.Item, u.Core.Soft.Value}, "")
	coreHard := strings.Join([]string{u.Core.Hard.Domain, u.Core.Hard.Type, u.Core.Hard.Item, u.Core.Hard.Value}, "")
	nofileSoft := strings.Join([]string{u.Nofile.Soft.Domain, u.Nofile.Soft.Type, u.Nofile.Soft.Item, u.Nofile.Soft.Value}, "")
	nofileHard := strings.Join([]string{u.Nofile.Hard.Domain, u.Nofile.Hard.Type, u.Nofile.Hard.Item, u.Nofile.Hard.Value}, "")
	nprocSoft := strings.Join([]string{u.Nproc.Soft.Domain, u.Nproc.Soft.Type, u.Nproc.Soft.Item, u.Nproc.Soft.Value}, "")
	nprocHard := strings.Join([]string{u.Nproc.Hard.Domain, u.Nproc.Hard.Type, u.Nproc.Hard.Item, u.Nproc.Hard.Value}, "")
	confLimit = append(confLimit, coreSoft, coreHard, nofileSoft, nofileHard, nprocSoft, nprocHard)
	s := ulimitCheck()
	if u.Core.State == 1 {
		compareUlimit(coreSoft, s)
		compareUlimit(coreHard, s)
	}
	if u.Nofile.State == 1 {
		compareUlimit(nofileSoft, s)
		compareUlimit(nofileHard, s)
	}
	if u.Nproc.State == 1 {
		compareUlimit(nprocSoft, s)
		compareUlimit(nprocHard, s)

	}
}

//dirModeCheck 检查目录或文件是否存在
func dirModeCheck(p string) (os.FileInfo, error) {
	fi, err := os.Lstat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	return fi, nil
}

// ulimitCheck 查询ulimit信息
func ulimitCheck() []string {
	var serverLimit []string
	f, err := os.Open("/etc/security/limits.conf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		s := strings.Join(strings.Fields(line), "")
		serverLimit = append(serverLimit, s)
	}
	return serverLimit
}

// compareUlimit 确认配置文件里的条目是否在server中存在
func compareUlimit(c string, s []string) {
	serverLimitStr := strings.Join(s, "")
	if !strings.Contains(serverLimitStr, c) {
		fmt.Println(c)
	}
}

// sysKernelCMD 校验内核参数
func sysKernelCMD(s []util.Syskernel) {
	defer wg.Done()
	if s[0].Value == 1 {
		for _, info := range s {
			if info.Name == "state" {
				continue
			}
			v := strconv.Itoa(info.Value)
			kMap, err := sysctl.GetPattern(info.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if v == kMap[info.Name] {
				continue
			} else {
				fmt.Printf("Server %s value %s\n", info.Name, kMap[info.Name])
			}
		}
	}
}

// dnsCMD 校验DNS信息
func dnsCMD(d *util.Dns) {
	defer wg.Done()
	if d.State == 1 {
		dnsSlice, err := util.GetDNSInfo()
		if err != nil {
			fmt.Println(err)
			return
		}
		dnsStr := strings.Join(dnsSlice, ",")
		for _, dns := range d.NameServer {
			if strings.Contains(dnsStr, dns) {
				continue
			}
			fmt.Printf("DNS MISSED %s", dns)
		}
	}
}
