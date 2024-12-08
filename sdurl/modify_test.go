package sdurl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	is := assert.New(t)

	url0 := "https://host1.com:3333/seg1/seg2?k2=v2&k1=v1"
	url1, err := Modify(url0)
	is.NoError(err)
	is.True(url0 == url1)

	// hostname
	url2, err := Modify(url0, SetHostname("host2.com"))
	is.NoError(err)
	url2a, err := url.Parse(url2)
	is.NoError(err)
	is.Equal("host2.com", url2a.Hostname())
	is.Equal("3333", url2a.Port())
	url2b, err := Modify(url2, SetHostname("host1.com"))
	is.NoError(err)
	is.True(url2b == url0)

	// port
	url3, err := Modify(url0, SetPort("4444"))
	is.NoError(err)
	url3a, err := url.Parse(url3)
	is.NoError(err)
	is.Equal("host1.com", url3a.Hostname())
	is.Equal("4444", url3a.Port())
	url3b, err := Modify(url3, SetPort("3333"))
	is.NoError(err)
	is.True(url3b == url0)

	// path
	url4, err := Modify(url0, SetPath("/abc"))
	is.NoError(err)
	url4a, err := url.Parse(url4)
	is.NoError(err)
	is.Equal("/abc", url4a.Path)
	url4b, err := Modify(url4, SetPath("/seg1/seg2"))
	is.NoError(err)
	is.True(url4b == url0)

	// query
	url5, err := Modify(url0, SetQuery("k3", "v3"))
	is.NoError(err)
	url5a, err := url.Parse(url5)
	is.NoError(err)
	is.Equal("v3", url5a.Query().Get("k3"))
	url5, err = Modify(url0, SetQuery("k2", "v2a"))
	is.NoError(err)
	url5a, err = url.Parse(url5)
	is.NoError(err)
	is.Equal("v2a", url5a.Query().Get("k2"))
	url5, err = Modify(url0, SetQueries(map[string]string{"k2": "v2b", "k3": "v3b"}))
	is.NoError(err)
	url5a, err = url.Parse(url5)
	is.NoError(err)
	is.Equal("v2b", url5a.Query().Get("k2"))
	is.Equal("v3b", url5a.Query().Get("k3"))
}
