package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stackBasics(t *testing.T) {
	st := NewStack[string]()
	a := "A"
	b := "B"
	c := "C"
	st.Push(&a, &b, &c)
	assert.Equal(t, 3, st.Length())
	e3 := st.Pop()
	assert.Equal(t, 2, st.Length())
	e2 := st.Pop()
	e1 := st.Pop()
	e0 := st.Pop()
	assert.Equal(t, 0, st.Length())
	assert.Equal(t, &c, e3)
	assert.Equal(t, &b, e2)
	assert.Equal(t, &a, e1)
	assert.Nil(t, e0)
}

func Test_stackAsSlice(t *testing.T) {
	st := NewStack[string]()
	a := "A"
	b := "B"
	st.Push(&a, &b)
	expected := []*string{&a, &b}
	assert.Equal(t, expected, st.AsSlice())
}
