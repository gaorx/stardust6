package sdhttpua

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/mssola/useragent"
)

// UserAgent 描述一个浏览器的User-Agent
type UserAgent struct {
	Raw                  string
	Mozilla              string
	Platform             string
	OS                   string
	Loc                  string
	Bot                  bool
	Mobile               bool
	BrowserEngine        string
	BrowserEngineVersion string
	BrowserName          string
	BrowserVersion       string
}

// Parse 解析User-Agent
func Parse(s string, to *UserAgent) (*UserAgent, error) {
	ua0 := useragent.New(s)
	if ua0 == nil {
		return nil, sderr.Wrap(ErrParse)
	}
	if to == nil {
		to = &UserAgent{}
	}
	to.Raw = ua0.UA()
	to.Mozilla = ua0.Mozilla()
	to.Platform = ua0.Platform()
	to.OS = ua0.OS()
	to.Loc = ua0.Localization()
	to.Bot = ua0.Bot()
	to.Mobile = ua0.Mobile()
	to.BrowserEngine, to.BrowserEngineVersion = ua0.Engine()
	to.BrowserName, to.BrowserVersion = ua0.Browser()
	return to, nil
}

// Match 判断和断言相匹配，多个断言之间是AND关系，必须满足所有断言才返回true
func (ua *UserAgent) Match(predicates ...Predicate) bool {
	if ua == nil {
		return false
	}
	if len(predicates) <= 0 {
		return true
	}
	for _, p := range predicates {
		if !p(ua) {
			return false
		}
	}
	return true
}
