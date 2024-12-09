package sdbackoff

import (
	"testing"

	"github.com/gaorx/stardust6/sderr"
	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	is := assert.New(t)

	n := 0
	err := Retry(Zero(), func() error {
		n++
		return nil
	})
	is.NoError(err)
	is.Equal(1, n)

	n = 0
	once := false
	err = Retry(Zero(), func() error {
		n++
		if once {
			return nil
		}
		once = true
		return sderr.Newf("some error")
	})
	is.NoError(err)
	is.Equal(2, n)
}

func TestStop(t *testing.T) {
	is := assert.New(t)

	n := 0
	err := Retry(Stop(), func() error {
		n++
		return nil
	})
	is.NoError(err)
	is.Equal(1, n)

	n = 0
	err = Retry(Stop(), func() error {
		n++
		return sderr.Newf("some error")
	})
	is.Error(err)
	is.Equal(1, n)
}
