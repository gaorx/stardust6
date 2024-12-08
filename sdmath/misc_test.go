package sdmath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalize(t *testing.T) {
	is := assert.New(t)
	is.Equal(4.0,
		Normalize[float64](
			2.0,
			Interval[float64]{1.0, 3.0},
			Interval[float64]{2.0, 6.0},
		),
	)
	is.Panics(func() {
		Normalize[float64](
			2.0,
			Interval[float64]{3.0, 3.0},
			Interval[float64]{2.0, 6.0},
		)
	})
}
