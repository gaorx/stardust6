package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReplaceBetweenTwoLines(t *testing.T) {
	is := assert.New(t)

	trim := func(s string) string {
		return strings.TrimSpace(s)
	}

	var replaced string
	var err error

	replaced, err = ReplaceBetweenTwoLines(
		trim(`
HEADER
	// start
	CONTENT1
	// end
----
// start
CONTENT2
// end
FOOTER
`),
		StringContains("// start"),
		HasSuffix("// end"),
		Text("NEW_CONTENT"),
	)
	is.NoError(err)
	is.Equal(trim(`
HEADER
	// start
NEW_CONTENT
	// end
----
// start
NEW_CONTENT
// end
FOOTER	
`), trim(replaced))

	replaced, err = ReplaceBetweenTwoLines(
		trim(`
HEADER
	// start
	CONTENT1
	// end
----
// start
CONTENT2
FOOTER
`),
		StringContains("// start"),
		StringContains("// end"),
		Text("NEW_CONTENT"),
	)
	is.Error(err)
}
