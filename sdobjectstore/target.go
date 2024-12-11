package sdobjectstore

import (
	"fmt"
	"path/filepath"
)

// Target 存储后的目标结果
type Target struct {
	Typ            TargetType
	Prefix         string
	InternalPrefix string
	Path           string
}

type TargetType int

const (
	DiscardTarget = TargetType(1)
	FileTarget    = TargetType(2)
	HttpTarget    = TargetType(3)
)

// Url 获取目标的URL
func (t *Target) Url() string {
	return t.url(false, false)
}

// HttpsUrl 获取目标的HTTPS URL
func (t *Target) HttpsUrl() string {
	return t.url(true, false)
}

// InternalUrl 获取目标的内部URL
func (t *Target) InternalUrl() string {
	return t.url(false, true)
}

func (t *Target) url(https bool, internal bool) string {
	switch t.Typ {
	case DiscardTarget:
		return ""
	case FileTarget:
		if https {
			return ""
		}
		absFn, err := filepath.Abs(filepath.Join(t.Prefix, t.Path))
		if err != nil {
			return ""
		}
		return "file://" + absFn
	case HttpTarget:
		var protocol, host string
		if https {
			protocol = "https"
		} else {
			protocol = "http"
		}
		if internal {
			host = t.InternalPrefix
		} else {
			host = t.Prefix
		}
		return fmt.Sprintf("%s://%s/%s", protocol, host, t.Path)
	default:
		return ""
	}
}
