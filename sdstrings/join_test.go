package sdstrings

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapJoin(t *testing.T) {
	is := assert.New(t)
	a := lo.Range(3)
	is.Equal("x1;x2;x3", JoinFunc(a, ";", func(v int, _ int) string {
		return fmt.Sprintf("x%d", v+1)
	}))
}
