/*
cpu信息，只包含cpu型号；cpu逻辑核数
*/
package util

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

// CPU 结构体，
type CpuInfo struct {
	ModelName string //cpu型号
	Processor int    //逻辑cpu id
}

// 获取CPU信息
func (cpu *CpuInfo) GetCpu() {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cpu.Processor = runtime.NumCPU()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		/*
			if strings.HasPrefix(line, "processor") {
				cpu.Processor = append(cpu.Processor, strings.TrimSpace(strings.Split(line, ":")[1]))

			}
		*/
		if strings.HasPrefix(line, "model name") {
			cpu.ModelName = strings.TrimSpace(strings.Split(line, ":")[1])

		}
	}
}
