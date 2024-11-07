package sdrand

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	is := assert.New(t)
	c := newCounter[string]()
	for i := 0; i < 1000000; i++ {
		l := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
		Shuffle(l)
		joined := strings.Join(l[0:3], "-")
		c.inc(joined)
	}
	is.Equal(7*8*9, c.size())
}
