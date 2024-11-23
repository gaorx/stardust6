package sdlocal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	is := assert.New(t)

	// IP4
	ips, err := IPs(Is4())
	is.NoError(err)
	for _, ip := range ips {
		is.True(len(ip.To4()) > 0)
	}

	// !IP4
	ips, err = IPs(Is4().Not())
	is.NoError(err)
	for _, ip := range ips {
		is.True(len(ip.To4()) <= 0)
	}

	// Loopback
	ips, err = IPs(IsLoopback())
	is.NoError(err)
	for _, ip := range ips {
		is.True(ip.IsLoopback())
	}

	// Private
	ips, err = IPs(IsPrivate())
	is.NoError(err)
	for _, ip := range ips {
		is.True(ip.IsPrivate())
	}
}
