//go:build !windows

package sdexec

import (
	"syscall"
)

// Rusage 程序运行所耗用的资源
type Rusage syscall.Rusage

func getRusage(ru any) Rusage {
	if ru1 := ru.(*syscall.Rusage); ru != nil {
		rur := *ru1
		return Rusage(rur)
	} else {
		return Rusage{}
	}
}
