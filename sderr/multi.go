package sderr

import (
	"fmt"
	"strings"
)

// MultiError 一种error，包含多个error
type MultiError struct {
	Errs []error
}

var _ error = (*MultiError)(nil)

// Join 将多个error合并成一个error
func Join(errs ...error) error {
	return Append(nil, errs...)
}

// Append 将多个error追加到一个error后面，并返回新的error
func Append(err error, errs ...error) error {
	var all []error
	add := func(err0 error) {
		if err0 != nil {
			all = append(all, err0)
		}
	}
	add(err)
	for _, err0 := range errs {
		add(err0)
	}
	if len(all) <= 0 {
		return nil
	}
	if len(all) == 1 {
		return all[0]
	}
	head, tails := all[0], all[1:]
	if head1, ok := ProbeMulti(head); ok {
		for _, err0 := range tails {
			if err1, ok := ProbeMulti(err0); ok {
				head1.Errs = append(head1.Errs, err1.Errs...)
			} else {
				head1.Errs = append(head1.Errs, err0)
			}
		}
		return head1
	} else {
		return &MultiError{Errs: all}
	}
}

// Error 实现Error.Error
func (e *MultiError) Error() string {
	if e == nil {
		return ""
	}
	lines := e.Lines("")
	switch len(lines) {
	case 0:
		return ""
	case 1:
		return lines[0]
	default:
		return fmt.Sprintf("(%s)", strings.Join(lines, " & "))
	}
}

// Lines 返回每个error的可读消息，每行一个，并且每行可以附加一个前缀
func (e *MultiError) Lines(prefix string) []string {
	if e == nil {
		return nil
	}
	errs := e.Errs
	var lines []string
	for _, err0 := range errs {
		if err0 != nil {
			msg := err0.Error()
			if msg != "" {
				lines = append(lines, prefix+msg)
			}
		}
	}
	return lines
}

// Unwrap 实现errors包的Unwrap
func (e *MultiError) Unwrap() []error {
	if e == nil {
		return nil
	}
	return e.Errs
}

// Empty 返回是否包含有其他错误信息
func (e *MultiError) Empty() bool {
	if e == nil {
		return true
	}
	return len(e.Errs) <= 0
}
