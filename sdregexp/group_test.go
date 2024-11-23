package sdregexp

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGroup(t *testing.T) {
	is := assert.New(t)
	r := regexp.MustCompile(`^(?P<first>\w+)-(?P<last>\w+)$`)
	group := FindStringSubmatchGroup(r, "hello-world")
	is.EqualValues(map[string]string{"first": "hello", "last": "world"}, group)
}
