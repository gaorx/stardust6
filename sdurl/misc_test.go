package sdurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitHostPort(t *testing.T) {
	is := assert.New(t)
	host, port := SplitHostPort("localhost:3390")
	is.Equal("localhost", host)
	is.Equal("3390", port)
	host, port = SplitHostPort("localhost")
	is.Equal("localhost", host)
	is.Equal("", port)
}

func TestEnsureHttpSchema(t *testing.T) {
	is := assert.New(t)
	is.Equal("http://localhost", EnsureHttpSchema("localhost", "http"))
	is.Equal("https://localhost", EnsureHttpSchema("localhost", "https"))
	is.Equal("http://localhost", EnsureHttpSchema("http://localhost", "http"))
	is.Equal("https://localhost", EnsureHttpSchema("https://localhost", "http"))
	is.Equal("http://localhost", EnsureHttpSchema("http://localhost", "https"))
	is.Equal("https://localhost", EnsureHttpSchema("https://localhost", "http"))
}
