package sdparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	is := assert.New(t)
	is.Equal(33, IntOr("33", 44))
	is.Equal(44, IntOr("error33", 44))
}
