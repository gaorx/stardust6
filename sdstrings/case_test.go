package sdstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	is := assert.New(t)
	is.Equal("hello_world", ToSnakeL("HelloWorld"))
	is.Equal("HELLO_WORLD", ToSnakeU("HelloWorld"))
	is.Equal("HELLO_WORLD", ToSnakeU("hello_world"))

	is.Equal("hello-world", ToKebabL("HelloWorld"))
	is.Equal("HELLO-WORLD", ToKebabU("HelloWorld"))
	is.Equal("HELLO-WORLD", ToKebabU("hello_world"))
	is.Equal("HELLO-WORLD", ToKebabU("hello-world"))

	is.Equal("hello-world", ToDelimitedL("HelloWorld", '-'))
	is.Equal("HELLO-WORLD", ToDelimitedU("Hello_World", '-'))

	is.Equal("helloWorld", ToCamelL("hello_world"))
	is.Equal("HelloWorld", ToCamelU("hello-world"))
	is.Equal("HelloWorld", ToCamelU("helloWorld"))

	is.EqualValues([]string{"Hello", "World"}, SplitCamel("HelloWorld"))
	is.EqualValues([]string{"hello", "World"}, SplitCamel("helloWorld"))
	is.EqualValues([]string{"RGB", "900"}, SplitCamel("RGB900"))
	is.EqualValues([]string{"hello", "_", "world"}, SplitCamel("hello_world"))
}
