package sdsync

import (
	"sync"
)

// Lock 锁定mtx，执行action，如果mtx为nil，则不锁定执行
func Lock(mtx *sync.Mutex, action func()) {
	if mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}
	action()
}

// LockW 锁定mtx的写操作，执行action，如果mtx为nil，则不锁定执行
func LockW(mtx *sync.RWMutex, action func()) {
	if mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}
	action()
}

// LockR 锁定mtx的读操作，执行action，如果mtx为nil，则不锁定执行
func LockR(mtx *sync.RWMutex, action func()) {
	if mtx != nil {
		mtx.RLock()
		defer mtx.RUnlock()
	}
	action()
}
