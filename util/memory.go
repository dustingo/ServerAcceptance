/*
内存信息
*/
package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// 内存单位KB
func GetMemory() (int, int, int, error) {
	var total, avail, swap string
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "MemTotal") {
			total = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "MemAvailable") {
			avail = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "SwapCached") {
			swap = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return trimKB(total), trimKB(avail), trimKB(swap), nil
}

func trimKB(memstr string) int {
	n, _ := strconv.Atoi(strings.Split(memstr, " ")[0])
	return n
}
