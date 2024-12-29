package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	is := assert.New(t)
	g := New()
	g.Add("a.txt", Text("A"))
	g.Sub("sub1").Also(func(sub Interface) {
		sub.Add("b.txt", Text("B"))
	})
	g.Sub("sub2").Also(func(sub Interface) {
		sub.Add("c.txt", Text("C"))
	})
	outs, err := g.TryAll(nil)
	is.NoError(err)
	is.Len(outs, 3)
	is.Equal("A", outs.Get("a.txt").Text())
	is.Equal("B", outs.Get("sub1/b.txt").Text())
	is.Equal("C", outs.Get("sub2/c.txt").Text())
}
