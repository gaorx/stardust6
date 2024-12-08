package sdsemver

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdlo"
	"github.com/gaorx/stardust6/sdparse"
	"github.com/gaorx/stardust6/sdstrings"
)

// V 语义化版本，格式为 major.minor.patch
type V struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

const numLimit = 10000

// New 创建一个版本
func New(major, minor, patch int) V {
	return V{Major: major, Minor: minor, Patch: patch}
}

// String 返回版本字符串
func (v V) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// IsZero 是否为空版本
func (v V) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

// IsValidate 是否合法
func (v V) IsValidate() bool {
	return (v.Major >= 0 && v.Major < numLimit) &&
		(v.Minor >= 0 && v.Minor < numLimit) &&
		(v.Patch >= 0 && v.Patch < numLimit)
}

// ToInt 转换为int64形式整数，便于存储在数据库中
func (v V) ToInt() int64 {
	return int64((v.Major * numLimit * numLimit) + (v.Minor * numLimit) + v.Patch)
}

// Equal 判断版本是是否是major.minor.patch
func (v V) Equal(major, minor, patch int) bool {
	return v.Major == major && v.Minor == minor && v.Patch == patch
}

// Parse 解析"x.y.z"形式的字符串到语义化版本
func Parse(s string) (V, error) {
	if s == "" {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver parse empty")
	}
	majorStr, minorStr, patchStr := sdstrings.Split3s(s, ".")
	major, err := sdparse.IntE(sdlo.EmptyOr(majorStr, "0"))
	if err != nil {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver parse major error")
	}
	minor, err := sdparse.IntE(sdlo.EmptyOr(minorStr, "0"))
	if err != nil {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver parse minor error")
	}
	patch, err := sdparse.IntE(sdlo.EmptyOr(patchStr, "0"))
	if err != nil {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver parse patch error")
	}
	if !(major >= 0 && major < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal major")
	}
	if !(minor >= 0 && minor < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal minor")
	}
	if !(patch >= 0 && patch < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal patch")
	}
	return V{Major: major, Minor: minor, Patch: patch}, nil
}

// FromInt 从int64整数形式转换为语义化版本
func FromInt(i int64) (V, error) {
	major := i / (numLimit * numLimit)
	minor := (i - major*numLimit*numLimit) / numLimit
	patch := i - major*numLimit*numLimit - minor*numLimit
	if !(major >= 0 && major < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal major")
	}
	if !(minor >= 0 && minor < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal minor")
	}
	if !(patch >= 0 && patch < numLimit) {
		return V{}, sderr.Wrapf(ErrParse, "sdsemver illegal patch")
	}
	return V{Major: int(major), Minor: int(minor), Patch: int(patch)}, nil
}
