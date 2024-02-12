package stack_test

import (
	"testing"

	"github.com/schmizzel/go-graphics/pkg/internal/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("Empty Stack", func(t *testing.T) {
		s := stack.New[int]()
		s.Push(1)
		s.Push(2)

		assert.False(t, s.IsEmpty())

		v, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 2, v)

		v, ok = s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 1, v)

		_, ok = s.Pop()
		assert.False(t, ok)
		assert.True(t, s.IsEmpty())
	})

	t.Run("Filled Stack", func(t *testing.T) {
		s := stack.New(1, 2)
		s.Push(3)

		v, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 3, v)

		v, ok = s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 2, v)

		v, ok = s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 1, v)

		_, ok = s.Pop()
		assert.False(t, ok)
	})
}
