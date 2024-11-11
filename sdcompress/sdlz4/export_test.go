package sdlz4

import (
	"bytes"
	"testing"

	"github.com/gaorx/stardust6/sdrand"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	is := assert.New(t)
	data0 := []byte(sdrand.String(1303, sdrand.LowerCaseAlphanumericCharset))
	for _, level := range AllLevels {
		data1, err := Compress(data0, level)
		is.NoError(err)
		data2, err := Uncompress(data1)
		is.NoError(err)
		is.True(bytes.Equal(data0, data2))
	}
}
