/*
系统信息，系统版本、os版本号、内核版本、架构
*/
package util

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

//系统信息

type SystemInfo struct {
	OSName    string //os名称
	OSRelease string //os 版本号
	OSKernel  string // os 内核版本
	OSArch    string // os架构
}

//获取系统基本信息
func (system *SystemInfo) GetSystem() {
	fversion, err := os.Open("/proc/version")
	if err != nil {
		panic(err)
	}
	frelease, err := os.Open("/etc/redhat-release")
	if err != nil {
		panic(err)
	}
	defer frelease.Close()
	defer fversion.Close()
	versionData, _ := ioutil.ReadAll(fversion)
	releaseData, _ := ioutil.ReadAll(frelease)
	verData := strings.Split(string(versionData), " ")
	relData := strings.Split(string(releaseData), " ")
	system.OSName = relData[0]
	system.OSRelease = relData[3]
	system.OSKernel = verData[2]
	system.OSArch = runtime.GOARCH

}
