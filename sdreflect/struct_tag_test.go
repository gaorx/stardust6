package sdreflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStructTag(t *testing.T) {
	is := assert.New(t)

	type AA struct {
		Name string `json:"name,omitempty" sd:"name"`
		Age  int    `json:"age" sd:"age"`
	}

	is.Equal("name", ParseStructFieldTag(reflect.TypeOf(AA{}), "Name", "sd"))
	is.Equal("age", ParseStructFieldTag(reflect.TypeOf(AA{}), "Age", "sd"))
	is.Equal("", ParseStructFieldTag(reflect.TypeOf(AA{}), "Name", "not_exists"))

	stags := ParseStructFieldTags(reflect.TypeOf(AA{}), "Name")
	is.NotNil(stags)

	// Get
	is.Equal("name,omitempty", stags.Get("json"))
	is.Equal("", stags.Get("json1"))

	// Keys
	is.EqualValues([]string{"json", "sd"}, stags.Keys())

	// Len
	is.Equal(2, stags.Len())

	// First
	is.Equal("name", stags.First("not_exists", "sd"))
	is.Equal("", stags.First("not_exists1", "not_exists2"))

	// Has
	is.True(stags.Has("sd") && stags.Has("json"))
	is.False(stags.Has("not_exists"))

	// HasOne
	is.False(stags.HasOne())
	is.False(stags.HasOne("not_exists"))
	is.True(stags.HasOne("not_exists", "sd"))

	// HasAll
	is.True(stags.HasAll())
	is.False(stags.HasAll("not_exists", "sd"))
	is.True(stags.HasAll("json", "sd"))
}
