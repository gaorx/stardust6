package sdparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	is := assert.New(t)
	t0 := time.Now()
	for _, layout := range timeLayoutsForParse {
		s0 := t0.Format(layout)
		t1, err := TimeE(s0)
		is.NoError(err)
		s1 := t1.Format(layout)
		is.True(s0 == s1)
	}
}
