package sderr

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

func TestIs(t *testing.T) {
	is := assert.New(t)
	var err error

	err = Newf("ERR1")
	is.True(Is(err, err))

	err = Newf("Error: %w", fs.ErrExist)
	is.True(Is(err, fs.ErrExist))

	err = Wrap(fs.ErrExist)
	is.True(Is(err, fs.ErrExist))

	err = Wrapf(fs.ErrExist, "Error: %w", assert.AnError)
	is.True(Is(err, fs.ErrExist))

	err = Join(fs.ErrExist, assert.AnError)
	is.True(Is(err, fs.ErrExist))
	err = Join(assert.AnError, fs.ErrExist)
	is.True(Is(err, fs.ErrExist))

	err = Recover(func() {
		panic(fs.ErrExist)
	})
	is.True(Is(err, fs.ErrExist))

	err = Recoverf(func() {
		panic(fs.ErrExist)
	}, "Error: %w", assert.AnError)
	is.True(Is(err, fs.ErrExist))
}

func TestAs(t *testing.T) {
	is := assert.New(t)
	var err error
	var anError error = &fs.PathError{Err: fs.ErrExist}

	asOk := func(err error) bool {
		return lo.T2(As[*fs.PathError](err)).B
	}

	err = Newf("error: %w", anError)
	is.True(asOk(err))

	err = Wrap(anError)
	is.True(asOk(err))

	err = Wrapf(anError, "Error: %w", assert.AnError)
	is.True(asOk(err))

	err = Join(anError, assert.AnError)
	is.True(asOk(err))

	err = Join(assert.AnError, anError)
	is.True(asOk(err))

	err = Recover(func() {
		panic(anError)
	})
	is.True(asOk(err))

	err = Recoverf(func() {
		panic(anError)
	}, "Error: %w", assert.AnError)
	is.True(asOk(err))
}

func TestEnsure(t *testing.T) {
	is := assert.New(t)
	var err error

	err = Ensure(nil)
	is.Nil(err)

	err = Ensure(fs.ErrExist)
	is.True(Is(err, fs.ErrExist))

	err = Ensure("Error1")
	is.Equal("Error1", err.Error())

	err = Ensure(1)
	is.Equal("1", err.Error())

	err = Ensure(reflect.ValueOf(assert.AnError))
	is.True(Is(err, assert.AnError))
}

func TestAttrs(t *testing.T) {
	is := assert.New(t)

	is.Equal(map[string]any(nil), Attrs(nil))

	err1 := With("k1", "v1").With("k0", "v0").Newf("ERR1")
	err2 := fmt.Errorf("ERR2: %w", err1)
	err3 := With("k2", "v2").Wrapf(err2, "ERR2")
	err4 := fmt.Errorf("ERR4:%w", err3)
	err5 := With("k3", "v3").With("k0", "v00").Wrapf(err4, "ERR3")
	err6 := fmt.Errorf("ERR6: %w", err5)
	is.Equal(map[string]any{
		"k0": "v00",
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}, Attrs(err6))
}

func TestAttr(t *testing.T) {
	is := assert.New(t)

	is.Equal(map[string]any(nil), Attrs(nil))

	err1 := With("k1", "v1").With("k0", "v0").Newf("ERR1")
	err2 := fmt.Errorf("ERR2: %w", err1)
	err3 := With("k2", "v2").Wrapf(err2, "ERR2")
	err4 := fmt.Errorf("ERR4:%w", err3)
	err5 := With("k3", "v3").With("k0", "v00").Wrapf(err4, "ERR3")
	err6 := fmt.Errorf("ERR6: %w", err5)

	is.Equal(lo.T2(any("v1"), true), lo.T2(Attr(err6, "k1")))
	is.Equal(lo.T2(any("v2"), true), lo.T2(Attr(err6, "k2")))
	is.Equal(lo.T2(any("v3"), true), lo.T2(Attr(err6, "k3")))
	is.Equal(lo.T2(any("v00"), true), lo.T2(Attr(err6, "k0")))
	is.Equal(lo.T2(any(nil), false), lo.T2(Attr(err6, "key_not_exist")))
}

func newErr1(err error) error {
	return Newf("ERR1: %w", err)
}

func TestRoot(t *testing.T) {
	is := assert.New(t)
	root := fs.ErrExist
	err1 := newErr1(root)
	err2 := With("k1", "v1").Wrapf(err1, "ERR2")
	err3 := fmt.Errorf("ERR3: %w", err2)
	err4 := Wrap(err3)
	is.True(Root(err4) == root)

	// root stack
	frames := RootStack(err4).Frames()
	is.True(len(frames) > 1)
	is.True(strings.Contains(frames[1].Func, "newErr1"))
}
