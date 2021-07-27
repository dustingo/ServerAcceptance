package work

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dustingo/ServerAcceptance/util"
)

const (
	lastURL = "https://raw.githubusercontent.com/dustingo/configtoml/main/lastcheck.toml"
)

var lconfig util.LastCheckInfo

func LastCheck() {
	err := Wget(lastURL, "lastcheck.toml")
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
	fmt.Println("Service.off:")
	serviceCMD(offState, offName, offLabel)
	fmt.Println("Service.on:")
	serviceCMD(onState, onName, onLabel)
	yumState, yumLabe, yumName, pipState, pipLabel, pipName, perlState, perlLabel, perlName := last.GetPackage()
	fmt.Println("Package.yum:")
	packageCMD(yumState, yumName, yumLabe)
	fmt.Println("Package.pip:")
	packageCMD(pipState, pipName, pipLabel)
	fmt.Println("Package.perl:")
	packageCMD(perlState, perlName, perlLabel)
}

//  serviceCMD 校验Service模块
func serviceCMD(status int, names []string, label string) {
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
							fmt.Printf("%s still running [ERR]\n", unit)
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
								fmt.Printf("%s not running [ERR]\n", unit)
							}
						}
					}
				} else {
					fmt.Printf("%s not running [ERR]\n", unit)
				}
			}
		}
	} else {
		fmt.Printf("Service.%s check passed\n", label)
	}
}

// packageCMD 校验Package 模块
func packageCMD(status int, names []string, label string) {
	if status == 1 {

	} else {
		fmt.Printf("Package.%s pass", label)
	}
}
