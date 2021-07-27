/*
GetDiskStats函数返回一个列表，列表内容为FilesystemStats结构体；
需要自己解析FilesystemStats结构题获取相关的数据
常用的参数为：
MountPoint挂载点；
Size 大小；
Free 空闲空间
Avail 可用空间
大小的单位为MB
*/
package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/sys/unix"
)

// 忽略的挂载点和文件系统类型
const (
	defIgnoredMountPoints = "^/(dev|proc|sys|var/lib/docker/.+)($|/)"
	defIgnoredFSTypes     = "^(autofs|binfmt_misc|rootfs|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|tmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
)

var (
	rootfsPath = "/"
)

// 标签结构体
type FilesystemLabels struct {
	Device     string
	MountPoint string
	FsType     string
}

// 总的文件系统结构体
type FilesystemStats struct {
	Labels            FilesystemLabels
	Size, Free, Avail float64
	Files, FilesFree  float64
}

// 获取挂载点详细信息，返回filesystemLabels切片，error
func mountPointDetails() ([]FilesystemLabels, error) {
	file, err := os.Open(("/proc/1/mounts"))
	if errors.Is(err, os.ErrNotExist) {
		// 两个挂在点,root 挂载点/proc/mounts 系统挂载点/proc/1/moounts
		fmt.Println("Reading root mouts failed,falling back to system mounts", err)
		file, err = os.Open(("/proc/mounts"))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseFilesystemLabels(file)
}

// 解析打开的文件系统内容返回filesystemLables
func parseFilesystemLabels(r io.Reader) ([]FilesystemLabels, error) {
	var filesystems []FilesystemLabels
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 4 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}
		parts[1] = strings.Replace(parts[1], "\\040", " ", -1)
		parts[1] = strings.Replace(parts[1], "\\011", "\t", -1)

		filesystems = append(filesystems, FilesystemLabels{
			Device:     parts[0],
			MountPoint: rootfsStripPrefix(parts[1]),
			FsType:     parts[2],
		})
	}
	return filesystems, scanner.Err()
}

// 实际抓取文件系统状态的函数，返回filesystemStats
func GetDiskStats() ([]FilesystemStats, error) {
	mps, err := mountPointDetails()
	if err != nil {
		fmt.Println(err.Error())
	}
	stats := []FilesystemStats{}
	for _, labels := range mps {
		if ok, _ := regexp.MatchString(defIgnoredMountPoints, labels.MountPoint); ok {
			continue
		}
		if ok, _ := regexp.MatchString(defIgnoredFSTypes, labels.FsType); ok {
			//fmt.Println("已忽略此类型", labels.fsType)
			continue
		}
		buf := new(unix.Statfs_t)
		err = unix.Statfs(rootfsFilePath(labels.MountPoint), buf)
		stats = append(stats, FilesystemStats{
			Labels:    labels,
			Size:      float64(buf.Blocks) * float64(buf.Bsize) / 1024 / 1024,
			Free:      float64(buf.Bfree) * float64(buf.Bsize) / 1024 / 1024,
			Avail:     float64(buf.Bavail) * float64(buf.Bsize) / 1024 / 1024,
			Files:     float64(buf.Files),
			FilesFree: float64(buf.Ffree),
		})
	}
	return stats, nil
}

// 功能函数
func rootfsStripPrefix(path string) string {
	if rootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, rootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
func rootfsFilePath(name string) string {
	return filepath.Join(rootfsPath, name)
}
