package sdvalidate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExport(t *testing.T) {
	is := assert.New(t)

	// Struct
	type testUser struct {
		Name string `validate:"required"`
		Age  int    `validate:"min=1,max=10"`
	}

	type testUser2 struct {
		Name string `validate:"tag-not-exists"`
	}

	err := Struct(nil)
	is.True(IsInvalidValidationErr(err))
	err = Struct("xx")
	is.True(IsInvalidValidationErr(err))
	is.Panics(func() {
		err = Struct(testUser2{})
	})
	err = Struct(testUser{})
	is.True(IsValidationErr(err))
	err = Struct(&testUser{})
	is.True(IsValidationErr(err))
	err = Struct(testUser{Name: "xx", Age: 9})
	is.NoError(err)
	err = Struct(testUser{Name: "", Age: 9})
	is.True(IsValidationErr(err))
	err = Struct(testUser{Name: "xx", Age: 12})
	is.True(IsValidationErr(err))
	err = Struct(&testUser{Name: "xx", Age: 9})
	is.NoError(err)
	err = Struct(&testUser{Name: "", Age: 9})
	is.True(IsValidationErr(err))
	err = Struct(&testUser{Name: "xx", Age: 12})
	is.True(IsValidationErr(err))

	// StructPartial
	err = StructPartial(&testUser{Name: "", Age: 20}, []string{"Name"})
	is.True(IsValidationErr(err))
	err = StructPartial(&testUser{Name: "xx", Age: 20}, []string{"Name"})
	is.NoError(err)
	err = Struct(&testUser{Name: "xx", Age: 20})
	is.True(IsValidationErr(err))

	// Var
	err = Var("", "required")
	is.True(IsValidationErr(err))
	err = Var("xx", "required")
	is.NoError(err)
	err = Var(0, "required")
	is.True(IsValidationErr(err))
	err = Var(1, "required")
	is.NoError(err)
}
