package sdjson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	is := assert.New(t)

	j1 := `{"a":1, "b": "b"}`
	j2 := `{"a":1, "b": "b"`

	type A struct {
		A int    `json:"a"`
		B string `json:"b"`
	}

	var a A
	var pa *A
	err := UnmarshalBytes([]byte(j1), &a)
	is.NoError(err)
	is.Equal(A{A: 1, B: "b"}, a)
	err = UnmarshalString(j1, &pa)
	is.NoError(err)
	is.Equal(A{A: 1, B: "b"}, *pa)
	err = UnmarshalBytes([]byte(j2), &a)
	is.Error(err)
	err = UnmarshalString(j2, &pa)
	is.Error(err)

	a1, err := UnmarshalBytesT[A]([]byte(j1))
	is.NoError(err)
	is.Equal(A{A: 1, B: "b"}, a1)
	a2, err := UnmarshalStringT[*A](j1)
	is.NoError(err)
	is.Equal(A{A: 1, B: "b"}, *a2)
	a3 := UnmarshalBytesDef[A]([]byte(j2), A{A: 3})
	is.Equal(A{A: 3}, a3)
	a4 := UnmarshalStringDef[*A](j2, &A{A: 4})
	is.Equal(A{A: 4}, *a4)

	a5, err := UnmarshalValueBytes([]byte(j1))
	is.NoError(err)
	is.Equal(Object{"a": 1.0, "b": "b"}, a5.AsObject())
	a6, err := UnmarshalValueString(j1)
	is.NoError(err)
	is.Equal(Object{"a": 1.0, "b": "b"}, a6.AsObject())
	a7 := UnmarshalValueBytesDef([]byte(j2), 3)
	is.Equal(V(3), a7)
	a8 := UnmarshalValueStringDef(j2, "xyz")
	is.Equal(V("xyz"), a8)
}
