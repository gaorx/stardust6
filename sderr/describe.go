package sderr

import (
	"encoding/json"
)

// DescribeOptions 描述一个错误的选项
type DescribeOptions struct {
	Unwrap         bool                 // 是否展开多层wrap
	Stack          bool                 // 是否加入stack信息
	FrameFormatter func(f Frame) string // 自定义如何格式化一个stack frame
}

// Description 对一个error的描述信息，其中包含root和多层wrap
type Description struct {
	Root  DescriptionItem   `json:"root"`
	Wraps []DescriptionItem `json:"wraps,omitempty"`
}

// DescriptionItem 描述wrap中的一层error信息
type DescriptionItem struct {
	Msg   string         `json:"msg"`
	Attrs map[string]any `json:"attrs,omitempty"`
	Stack []string       `json:"stack,omitempty"`
}

// Describe 描述一个error,生成一个结构化的信息
func Describe(err error, opts *DescribeOptions) *Description {
	if err == nil {
		return nil
	}

	opts = ensurePtr(opts)
	if opts.FrameFormatter == nil {
		opts.FrameFormatter = func(f Frame) string {
			return f.String()
		}
	}

	makeItem := func(msg string, attrs map[string]any, stack *Stack) DescriptionItem {
		var frameLines []string
		if opts.Stack {
			frames := stack.Frames()
			for _, frame := range frames {
				frameLines = append(frameLines, opts.FrameFormatter(frame))
			}
		}
		return DescriptionItem{
			Msg:   msg,
			Attrs: attrs,
			Stack: frameLines,
		}
	}

	var res Description
	if !opts.Unwrap {
		if e, ok := Probe(err); ok {
			res.Root = makeItem(e.Error(), Attrs(e), RootStack(err))
		} else {
			res.Root = makeItem(err.Error(), nil, nil)
		}
	} else {
		unwrappedErrs := UnwrapNested(err)
		for i := 0; i < len(unwrappedErrs); i++ {
			unwrappedErr := unwrappedErrs[i]
			isRoot := i >= len(unwrappedErrs)-1
			if e, ok := Probe(unwrappedErr); ok {
				if isRoot {
					res.Root = makeItem(e.Msg(), e.OwnAttrs(), e.stack)
				} else {
					res.Wraps = append(res.Wraps, makeItem(e.Msg(), e.OwnAttrs(), e.stack))
				}
			} else {
				if isRoot {
					res.Root = makeItem(unwrappedErr.Error(), nil, nil)
				} else {
					res.Wraps = append(res.Wraps, makeItem(unwrappedErr.Error(), nil, nil))
				}
			}
		}
	}
	return &res
}

// Json 将一个Description转换成JSON
func (d *Description) Json(pretty bool) string {
	if d == nil {
		return ""
	}
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(d, "", "  ")
	} else {
		b, err = json.Marshal(d)
	}
	if err != nil {
		panic("failed to marshal description to json: " + err.Error())
	}
	return string(b)
}
