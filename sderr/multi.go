package sderr

import (
	"fmt"
	"strings"
)

type MultiError struct {
	Errs []error
}

var _ error = (*MultiError)(nil)

func Join(errs ...error) error {
	return Append(nil, errs...)
}

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

func (e *MultiError) Unwrap() []error {
	if e == nil {
		return nil
	}
	return e.Errs
}

func (e *MultiError) Empty() bool {
	if e == nil {
		return true
	}
	return len(e.Errs) <= 0
}
