package sdbytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSummarize(t *testing.T) {
	is := assert.New(t)
	a := []byte("summarize")
	is.Equal("bytes(0) nil", Summarize(nil))
	is.Equal("bytes(0) []", Summarize([]byte{}))
	is.Equal("bytes(1) [73]", Summarize(a[:1]))
	is.Equal("bytes(2) [73 75]", Summarize(a[:2]))
	is.Equal("bytes(3) [73 75 6d]", Summarize(a[:3]))
	is.Equal("bytes(4) [73 75 6d 6d]", Summarize(a[:4]))
	is.Equal("bytes(5) [73 75 6d 6d 61]", Summarize(a[:5]))
	is.Equal("bytes(6) [73 75 6d 6d 61 72]", Summarize(a[:6]))
	is.Equal("bytes(9) [73 75 6d 6d .. 7a 65]", Summarize(a))
}
