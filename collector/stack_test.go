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

func Test_stackGetOutOfRange(t *testing.T) {
	st := NewStack[string]()
	a := "A"
	st.Push(&a)
	assert.Nil(t, st.Get(-1))
	assert.Nil(t, st.Get(1))
	assert.Equal(t, &a, st.Get(0))
}

func Test_stackAsSliceIsCopy(t *testing.T) {
	st := NewStack[string]()
	a, b, c := "A", "B", "C"
	st.Push(&a, &b)
	snapshot := st.AsSlice()
	st.Pop()
	st.Push(&c)
	assert.Equal(t, []*string{&a, &b}, snapshot)
	assert.Equal(t, []*string{&a, &c}, st.AsSlice())
}
