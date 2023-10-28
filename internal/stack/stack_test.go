package stack_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisguycodes/modelstack/internal/stack"
)

func TestPushPeekPop(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()
	i := rand.Int()
	s.Push(i)

	assert.Equal(t, i, s.Peek().Value)
	assert.Equal(t, i, s.Pop().Value)
}

func TestSwapPeek(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()
	i := rand.Int()
	j := rand.Int()

	s.Push(i)
	s.Swap(j)

	assert.Equal(t, j, s.Peek().Value)
}

func TestEmptyPop(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()

	var n *stack.Element[int]
	assert.Equal(t, n, s.Pop())
}

func TestEmptyPeek(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()

	var n *stack.Element[int]
	assert.Equal(t, n, s.Peek())
}

func TestPopPop(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()

	i := rand.Int()
	j := rand.Int()

	s.Push(i)
	s.Push(j)

	s.Pop()

	assert.Equal(t, i, s.Pop().Value)
}

func TestPopPopCopySafe(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()

	i := rand.Int()
	j := rand.Int()

	s.Push(i)
	s.Push(j)

	cop := s
	assert.Equal(t, j, cop.Pop().Value)
	assert.Equal(t, i, s.Pop().Value)
}

func TestPushPopCopySafe(t *testing.T) {
	t.Parallel()

	s := stack.New[int]()

	i := rand.Int()
	j := rand.Int()

	cop := s

	s.Push(i)
	cop.Push(j)

	assert.Equal(t, j, s.Pop().Value)
	assert.Equal(t, i, cop.Pop().Value)
}

func TestNewPreSeed(t *testing.T) {
	t.Parallel()

	i := rand.Int()
	j := rand.Int()
	k := rand.Int()

	s := stack.New(i, j, k)

	assert.Equal(t, i, s.Pop().Value)
	assert.Equal(t, j, s.Pop().Value)
	assert.Equal(t, k, s.Pop().Value)
}
