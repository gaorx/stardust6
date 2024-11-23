package sdlo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyOr(t *testing.T) {
	is := assert.New(t)
	is.Equal(1, EmptyOr(1, 2))
	is.Equal(2, EmptyOr(0, 2))
}
