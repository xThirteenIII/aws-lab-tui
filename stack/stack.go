package stack

// Stack is a LIFO data structure. items is a slice of generics data.
type Stack[T any] struct {
	items []T
}

// NewStack returns an empty Stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// IsEmpty returns true if the stack is empty.
func (s Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Peek returns the last element of the Stack.
// Returns default nil value of T type, if Stack is empty.
func (s Stack[T]) Peek() T {
	if len(s.items) == 0 {
		var zero T
		return zero
	}
	return s.items[len(s.items)-1]
}

// Pop removes the last element of the Stack and returns it.
// Returns default nil value of T type, if Stack is empty.
func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		var zero T
		return zero
	}
	lastElement := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return lastElement
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}
