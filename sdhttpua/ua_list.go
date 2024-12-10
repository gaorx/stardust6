package sdhttpua

import (
	"github.com/gaorx/stardust6/sdcsv"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/samber/lo"
	"slices"
	"strconv"
)

// UserAgents User-Agent列表
type UserAgents struct {
	list     []*UserAgent
	weighted []sdrand.W[*UserAgent]
}

// FromSlice 从UserAgent列表创建UserAgents
func FromSlice(l []*UserAgent) *UserAgents {
	if len(l) <= 0 {
		return &UserAgents{}
	}
	return &UserAgents{
		list:     slices.Clone(l),
		weighted: nil,
	}
}

// FromLines 从文本行创建UserAgents，每行一个User-Agent
func FromLines(lines []string, ignoreErr bool) (*UserAgents, error) {
	if len(lines) <= 0 {
		return &UserAgents{}, nil
	}
	var l []*UserAgent
	for _, line := range lines {
		ua, err := Parse(line, nil)
		if err != nil {
			if ignoreErr {
				continue
			} else {
				return nil, sderr.Wrap(err)
			}
		}
		l = append(l, ua)
	}
	return &UserAgents{list: l, weighted: nil}, nil
}

// ParseLines 解析文本行，每行一个User-Agent
func ParseLines(text string, ignoreErr bool) (*UserAgents, error) {
	lines := sdstrings.SplitNonempty(text, "\n", true)
	return FromLines(lines, ignoreErr)
}

// ParseCSV 从CSV文本创建UserAgents，每个记录包含两个字段：weight和raw
func ParseCSV(text string, ignoreErr bool) (*UserAgents, error) {
	if text == "" {
		return &UserAgents{}, nil
	}
	reader, err := sdcsv.NewReaderText(text, &sdcsv.Options{
		Header:           false,
		Fields:           []string{"weight", "raw"},
		TrimLeadingSpace: true,
	})
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	recordMaps, err := reader.ReadMaps()
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	var l []*UserAgent
	var wl []sdrand.W[*UserAgent]
	for _, m := range recordMaps {
		ua, parseErr := Parse(m["raw"], nil)
		if parseErr != nil {
			if ignoreErr {
				continue
			} else {
				return nil, sderr.Wrap(err)
			}
		}
		weight, weightErr := strconv.Atoi(m["weight"])
		if weightErr != nil {
			weight = 0
		}
		l = append(l, ua)
		wl = append(wl, sdrand.W[*UserAgent]{W: weight, V: ua})
	}
	hasWeight := lo.ContainsBy(wl, func(w sdrand.W[*UserAgent]) bool { return w.W > 0 })
	if !hasWeight {
		wl = nil
	}
	return &UserAgents{list: l, weighted: wl}, nil
}

// IsEmpty 是否为空
func (ual *UserAgents) IsEmpty() bool {
	return ual == nil || len(ual.list) <= 0
}

// Len 返回UserAgent数量
func (ual *UserAgents) Len() int {
	if ual == nil {
		return 0
	}
	return len(ual.list)
}

// All 返回所有UserAgent
func (ual *UserAgents) All() []*UserAgent {
	if ual == nil {
		return nil
	}
	return slices.Clone(ual.list)
}

// At 返回指定索引的UserAgent
func (ual *UserAgents) At(index int) *UserAgent {
	if ual.IsEmpty() {
		return nil
	}
	if index < 0 || index >= len(ual.list) {
		return nil
	}
	return ual.list[index]
}

// Sub 返回满足断言的UserAgent列表作为子列表
func (ual *UserAgents) Sub(predicates ...Predicate) *UserAgents {
	l := ual.Find(predicates...)
	if len(l) <= 0 {
		return &UserAgents{}
	}
	var wl []sdrand.W[*UserAgent]
	if len(ual.weighted) > 0 {
		m := lo.SliceToMap(ual.weighted, func(w sdrand.W[*UserAgent]) (string, sdrand.W[*UserAgent]) {
			return w.V.Raw, w
		})
		for _, ua := range l {
			if uaw, ok := m[ua.Raw]; ok && uaw.W > 0 {
				wl = append(wl, uaw)
			}
		}
	}
	return &UserAgents{list: l, weighted: wl}
}

// Find 查找满足断言的UserAgent列表
func (ual *UserAgents) Find(predicates ...Predicate) []*UserAgent {
	if ual.IsEmpty() {
		return nil
	}
	var l []*UserAgent
	for _, ua := range ual.list {
		if ua.Match(predicates...) {
			l = append(l, ua)
		}
	}
	return l
}

// FindRaw 查找满足断言的User-Agent原始字符串列表
func (ual *UserAgents) FindRaw(predicates ...Predicate) []string {
	l := ual.Find(predicates...)
	return lo.Map(l, func(ua *UserAgent, _ int) string { return ua.Raw })
}

// Sample 随机抽样一个UserAgent，如果有权重信息，则按权重抽样，否则是随机抽样
func (ual *UserAgents) Sample() *UserAgent {
	if ual.IsEmpty() {
		return nil
	}
	if len(ual.weighted) > 0 {
		return sdrand.SampleWeighted[*UserAgent](ual.weighted...)
	} else {
		return sdrand.Sample(ual.list...)
	}
}

// SampleRaw 随机抽样一个User-Agent原始字符串，如果有权重信息，则按权重抽样，否则是随机抽样
func (ual *UserAgents) SampleRaw() string {
	ua := ual.Sample()
	if ua == nil {
		return ""
	}
	return ua.Raw
}

// Samples 随机抽样多个UserAgent，如果有权重信息，则按权重抽样，否则是随机抽样
func (ual *UserAgents) Samples(n int) []*UserAgent {
	var l []*UserAgent
	for i := 0; i < n; i++ {
		ua := ual.Sample()
		if ua != nil {
			l = append(l, ua)
		}
	}
	return l
}

// SamplesRaw 随机抽样多个User-Agent原始字符串，如果有权重信息，则按权重抽样，否则是随机抽样
func (ual *UserAgents) SamplesRaw(n int) []string {
	l := ual.Samples(n)
	return lo.Map(l, func(ua *UserAgent, _ int) string { return ua.Raw })
}

// Platforms 返回这些User-Agent列表中的所有Platform值
func (ual *UserAgents) Platforms() []string {
	if ual.IsEmpty() {
		return nil
	}
	return uniqAndSort(lo.Map(ual.list, func(ua *UserAgent, _ int) string { return ua.Platform }))
}

// OSes 返回这些User-Agent列表中的所有OS值
func (ual *UserAgents) OSes() []string {
	if ual.IsEmpty() {
		return nil
	}
	return uniqAndSort(lo.Map(ual.list, func(ua *UserAgent, _ int) string { return ua.OS }))
}

// BrowserEngines 返回这些User-Agent列表中的所有BrowserEngine值
func (ual *UserAgents) BrowserEngines() []string {
	if ual.IsEmpty() {
		return nil
	}
	return uniqAndSort(lo.Map(ual.list, func(ua *UserAgent, _ int) string { return ua.BrowserEngine }))
}

// BrowserNames 返回这些User-Agent列表中的所有BrowserName值
func (ual *UserAgents) BrowserNames() []string {
	if ual.IsEmpty() {
		return nil
	}
	return uniqAndSort(lo.Map(ual.list, func(ua *UserAgent, _ int) string { return ua.BrowserName }))
}
