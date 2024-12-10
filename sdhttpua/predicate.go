package sdhttpua

// Predicate 筛选UA用的断言
type Predicate func(*UserAgent) bool

// Not 取反的断言
func (p Predicate) Not() Predicate {
	return func(ua *UserAgent) bool {
		return !p(ua)
	}
}

// Or 逻辑OR的断言
func Or(first Predicate, others ...Predicate) Predicate {
	if len(others) <= 0 {
		return first
	}
	return func(ua *UserAgent) bool {
		if first != nil && first(ua) {
			return true
		}
		for _, other := range others {
			if other != nil && other(ua) {
				return true
			}
		}
		return false
	}
}

// PlatformIs 判断平台是某个值的断言，可以指定多个值，只要有一个匹配即可
func PlatformIs(platform string, others ...string) Predicate {
	return func(ua *UserAgent) bool {
		if ua.Platform == platform {
			return true
		}
		for _, other := range others {
			if other != "" && ua.Platform == other {
				return true
			}
		}
		return false
	}
}

// PlatformIsWindows 判断平台是Windows的断言
func PlatformIsWindows() Predicate {
	return PlatformIs("Windows")
}

// PlatformIsMacintosh 判断平台是Macintosh的断言
func PlatformIsMacintosh() Predicate {
	return PlatformIs("Macintosh")
}

// PlatformIsLinux 判断平台是Linux的断言
func PlatformIsLinux() Predicate {
	return PlatformIs("Linux")
}

// OSIs 判断OS是某个值的断言，可以指定多个值，只要有一个匹配即可
func OSIs(os string, others ...string) Predicate {
	return func(ua *UserAgent) bool {
		if ua.OS == os {
			return true
		}
		for _, other := range others {
			if other != "" && ua.OS == other {
				return true
			}
		}
		return false
	}
}

// BrowserNameIs 判断浏览器是某个值的断言，可以指定多个值，只要有一个匹配即可
func BrowserNameIs(name string, others ...string) Predicate {
	return func(ua *UserAgent) bool {
		if ua.BrowserName == name {
			return true
		}
		for _, other := range others {
			if other != "" && ua.BrowserName == other {
				return true
			}
		}
		return false
	}
}

// BrowserNameIsChrome 判断浏览器是Chrome的断言
func BrowserNameIsChrome() Predicate {
	return BrowserNameIs("Chrome")
}

// BrowserNameIsEdge 判断浏览器是Edge的断言
func BrowserNameIsEdge() Predicate {
	return BrowserNameIs("Edge")
}

// BrowserNameIsFirefox 判断浏览器是Firefox的断言
func BrowserNameIsFirefox() Predicate {
	return BrowserNameIs("Firefox")
}

// BrowserNameIsSafari 判断浏览器是Safari的断言
func BrowserNameIsSafari() Predicate {
	return BrowserNameIs("Safari")
}

// IsMobile 判断是否是移动端的断言
func IsMobile() Predicate {
	return func(ua *UserAgent) bool {
		return ua.Mobile
	}
}
