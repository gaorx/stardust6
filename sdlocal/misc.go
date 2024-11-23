package sdlocal

import (
	"os"
	"os/user"
	"runtime"
)

// Hostname 获取主机名
func Hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hn
}

// HomeDir 获取当前用户的家目录
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir
}

// OS 获取操作系统
func OS() string {
	return runtime.GOOS
}

// Arch 获取架构
func Arch() string {
	return runtime.GOARCH
}

// NumCPU 获取CPU数量
func NumCPU() int {
	return runtime.NumCPU()
}
