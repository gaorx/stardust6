package sdurl

import (
	"github.com/gaorx/stardust6/sderr"
	"net/url"
)

// Modifier URL修改器
type Modifier func(u *url.URL)

// Modify 通过一组修改器依次修改URL
func Modify(rawUrl string, modifiers ...Modifier) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", sderr.With("url", rawUrl).Wrapf(err, "parse url for modify error")
	}
	for _, f := range modifiers {
		if f != nil {
			f(u)
		}
	}
	return u.String(), nil
}

// ModifyOr 通过一组修改器依次修改URL，如果出错则返回默认值
func ModifyOr(rawUrl string, def string, modifiers ...Modifier) string {
	r, err := Modify(rawUrl, modifiers...)
	if err != nil {
		return def
	}
	return r
}

// SetQuery 修改器，设置查询中的一个参数
func SetQuery(k, v string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		q.Set(k, v)
		u.RawQuery = q.Encode()
	}
}

// SetQueries 修改器，设置查询中的多个参数
func SetQueries(m map[string]string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		for k, v := range m {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
}

// DeleteQuery 修改器，添加查询中的一个参数
func DeleteQuery(k string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		q.Del(k)
		u.RawQuery = q.Encode()
	}
}

// DeleteQueries 修改器，删除查询中的多个参数
func DeleteQueries(keys ...string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		for _, k := range keys {
			q.Del(k)
		}
		u.RawQuery = q.Encode()
	}
}

// SetPath 修改器，设置Path
func SetPath(path string) Modifier {
	return func(u *url.URL) {
		u.Path = path
	}
}

// SetHost 修改器，设置Host
func SetHost(host string) Modifier {
	return func(u *url.URL) {
		u.Host = host
	}
}

// SetHostname 修改器，设置Hostname
func SetHostname(hostname string) Modifier {
	return func(u *url.URL) {
		_, port := SplitHostPort(u.Host)
		if port != "" {
			u.Host = hostname + ":" + port
		} else {
			u.Host = hostname
		}
	}
}

// SetPort 修改器，设置Port
func SetPort(port string) Modifier {
	return func(u *url.URL) {
		hostname, _ := SplitHostPort(u.Host)
		if port != "" {
			u.Host = hostname + ":" + port
		} else {
			u.Host = hostname
		}
	}
}
