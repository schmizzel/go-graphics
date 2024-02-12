package stack

// Simple and efficient stack implementation
// Note that this implementation is not concurrency safe
type Stack[T any] struct {
	items []T
	head  int
}

// Creates a new stack with the given itemsack
func New[T any](items ...T) *Stack[T] {
	return &Stack[T]{
		items: items,
		head:  len(items) - 1,
	}
}

// Pushes the given item onto the stack
func (s *Stack[T]) Push(v T) {
	s.head++
	if s.head >= len(s.items) {
		s.items = append(s.items, v)
		return
	}

	s.items[s.head] = v
}

// Pops the top item from the stack and returns it
// Returns false if the stack is empty and true otherwise
func (s *Stack[T]) Pop() (T, bool) {
	if s.head < 0 {
		var zero T
		return zero, false
	}

	s.head--
	return s.items[s.head+1], true
}

// Returns true if the stack is empty and false otherwise.
func (s *Stack[T]) IsEmpty() bool {
	return s.head < 0
}
