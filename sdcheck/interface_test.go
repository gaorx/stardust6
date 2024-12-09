package sdcheck

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAll(t *testing.T) {
	is := assert.New(t)

	// func
	is.NoError(Func(nil).Check())
	is.NoError(Func(func() error {
		return nil
	}).Check())
	is.Error(Func(func() error {
		return sderr.Newf("FUNC")
	}).Check())

	// true
	is.NoError(True(true, "TRUE").Check())
	is.Error(True(false, "TRUE").Check())

	// false
	is.NoError(False(false, "FALSE").Check())
	is.Error(False(true, "FALSE").Check())

	// not
	is.Error(Not(True(true, "TRUE"), "NOT").Check())
	is.NoError(Not(True(false, "TRUE"), "NOT").Check())

	// all
	is.NoError(All().Check())
	is.NoError(All(
		True(true, "TRUE"),
		True(true, "TRUE"),
	).Check())
	is.Error(All(
		True(false, "TRUE"),
		True(true, "TRUE"),
	).Check())
	is.Error(All(
		True(true, "TRUE"),
		True(false, "TRUE"),
	).Check())

	// And
	is.NoError(And(nil, "AND").Check())
	is.NoError(And([]Interface{
		True(true, "TRUE"),
	}, "AND").Check(), "AND")
	is.Error(And([]Interface{
		True(false, "TRUE"),
	}, "AND").Check(), "AND")
	is.NoError(And([]Interface{
		True(true, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	is.Error(And([]Interface{
		True(false, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	is.Error(And([]Interface{
		True(true, "TRUE"),
		True(false, "TRUE"),
	}, "AND").Check())

	// Or
	is.NoError(Or(nil, "OR").Check())
	is.NoError(Or([]Interface{
		True(true, "TRUE"),
	}, "OR").Check())
	is.Error(Or([]Interface{
		True(false, "TRUE"),
	}, "OR").Check())
	is.NoError(Or([]Interface{
		True(true, "TRUE"),
		True(true, "TRUE"),
	}, "OR").Check())
	is.NoError(Or([]Interface{
		True(false, "TRUE"),
		True(true, "TRUE"),
	}, "AND").Check())
	is.NoError(Or([]Interface{
		True(true, "TRUE"),
		True(false, "TRUE"),
	}, "OR").Check())
	is.Error(Or([]Interface{
		True(false, "TRUE"),
		True(false, "TRUE"),
	}, "OR").Check())

	// if
	is.NoError(If(true, True(true, "TRUE")).Check())
	is.NoError(If(false, True(true, "TRUE")).Check())
	is.Error(If(true, True(false, "TRUE")).Check())
	is.NoError(If(false, True(false, "TRUE")).Check())

	// For
	var a, b int
	a, b = 0, 0
	is.NoError(All(
		For(func() (int, error) { return 3, nil }, &a),
		For(func() (int, error) { return 4, nil }, &b),
	).Check())
	is.Equal(3, a)
	is.Equal(4, b)
	a, b = 0, 0
	is.Error(All(
		For(func() (int, error) { return 3, sderr.Newf("xx") }, &a),
		For(func() (int, error) { return 4, nil }, &b),
	).Check())
	is.Equal(0, a)
	is.Equal(0, b)

	// lazy
	a, b = 0, 0
	is.NoError(All(
		For(func() (int, error) { return 3, nil }, &a),
		Lazy(func() Interface {
			// 在lazy中可以使用被上一个checker修改过的a，不放在lazy中则不行
			return For(func() (int, error) { return 7 + a, nil }, &b)
		}),
	).Check())
	is.Equal(3, a)
	is.Equal(10, b)
	a, b = 0, 0
	is.Error(All(
		For(func() (int, error) { return 3, nil }, &a),
		Lazy(func() Interface {
			return For(func() (int, error) { return 7 + a, sderr.Newf("XX") }, &b)
		}),
	).Check())
	is.Equal(a, 3)
	is.Equal(b, 0)

	// required
	is.Error(Required(nil, "REQUIRED").Check())
	is.Error(Required((func() int)(nil), "REQUIRED").Check())
	is.Error(Required((*int)(nil), "REQUIRED").Check())
	is.NoError(Required(1, "REQUIRED").Check())
	is.Error(Required(0, "REQUIRED").Check())
	is.NoError(Required(true, "REQUIRED").Check())
	is.Error(Required(false, "REQUIRED").Check())
	is.NoError(Required("a", "REQUIRED").Check())
	is.Error(Required("", "REQUIRED").Check())
	is.Error(Required([]int{}, "REQUIRED").Check())
	is.Error(Required([]int(nil), "REQUIRED").Check())
	is.NoError(Required([]int{0}, "REQUIRED").Check())
	is.Error(Required(map[string]int{}, "REQUIRED").Check())
	is.Error(Required(map[string]int(nil), "REQUIRED").Check())
	is.NoError(Required(map[string]int{"": 0}, "REQUIRED").Check())
	is.Error(Required(struct{}{}, "REQUIRED").Check())
	is.NoError(Required(&struct{}{}, "REQUIRED").Check())

	// len
	is.NoError(Len([]int{}, 0, 2, "LEN").Check())
	is.NoError(Len([]int{0}, 0, 2, "LEN").Check())
	is.NoError(Len([]int{0, 0}, 0, 2, "LEN").Check())
	is.Error(Len([]int{0, 0, 0}, 0, 2, "LEN").Check())
	is.NoError(Len("", 0, 2, "LEN").Check())
	is.NoError(Len("a", 0, 2, "LEN").Check())
	is.NoError(Len("aa", 0, 2, "LEN").Check())
	is.Error(Len("aaa", 0, 2, "LEN").Check())

	// collection
	is.Error(In("a", []string{}, "IN").Check())
	is.NoError(NotIn("a", []string{}, "NOT_IN").Check())
	is.NoError(In("a", []string{"b", "a"}, "IN").Check())
	is.Error(In("a", []string{"b", "c"}, "IN").Check())
	is.Error(NotIn("a", []string{"b", "a"}, "NOT_IN").Check())
	is.NoError(NotIn("a", []string{"b", "c"}, "NOT_IN").Check())
	is.NoError(HasKey("a", map[string]int{"a": 0}, "HAS_KEY").Check())
	is.Error(HasKey("a", map[string]int{"b": 0}, "HAS_KEY").Check())
	is.Error(NotHasKey("a", map[string]int{"a": 0}, "NOT_HAS_KEY").Check())
	is.NoError(NotHasKey("a", map[string]int{"b": 0}, "NOT_HAS_KEY").Check())

	// string
	is.NoError(MatchRegexp("abc", "a[bd]c", "MATCH_REGEXP").Check())
	is.NoError(MatchRegexp("adc", "a[bd]c", "MATCH_REGEXP").Check())
	is.Error(MatchRegexp("aec", "a[bd]c", "MATCH_REGEXP").Check())
	is.NoError(MatchRegexpPattern("abc", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	is.NoError(MatchRegexpPattern("adc", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	is.Error(MatchRegexpPattern("aec", regexp.MustCompile("a([bd])c"), "MATCH_REGEXP_PATTERN").Check())
	is.NoError(HasSub("abc", "b", "HAS_SUB").Check())
	is.Error(HasSub("adc", "b", "HAS_SUB").Check())
	is.NoError(HasPrefix("ab", "a", "HAS_PREFIX").Check())
	is.Error(HasPrefix("cb", "a", "HAS_PREFIX").Check())
	is.NoError(HasSuffix("ba", "a", "HAS_PREFIX").Check())
	is.Error(HasSuffix("bc", "a", "HAS_SUFFIX").Check())

	// Validate
	type testUser struct {
		Name string `validate:"required"`
		Age  int    `validate:"min=1,max=10"`
	}
	is.Error(ValidateStruct(nil, "VALIDATE").Check())
	is.Error(ValidateStruct(&testUser{Name: "", Age: 20}, "VALIDATE").Check())
	is.Error(ValidateStruct(&testUser{Name: "xx", Age: 20}, "VALIDATE").Check())
	is.NoError(ValidateStruct(&testUser{Name: "xx", Age: 9}, "VALIDATE").Check())
	is.Error(ValidateStructPartial(&testUser{Name: "", Age: 20}, []string{"Name"}, "VALIDATE").Check())
	is.NoError(ValidateStructPartial(&testUser{Name: "xx", Age: 20}, []string{"Name"}, "VALIDATE").Check())
	is.Error(ValidateVar(20, "min=1,max=10", "VALIDATE").Check())
	is.NoError(ValidateVar(8, "min=1,max=10", "VALIDATE").Check())
}
