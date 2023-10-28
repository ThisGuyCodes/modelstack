package stack

// Element is an Element of a stack.
type Element[T any] struct {
	// Next Element in the stack
	next *Element[T]

	// The Value stored with this element.
	Value T
}

func New[T any](init ...T) *Stack[T] {
	s := &Stack[T]{}

	if len(init) > 0 {
		elems := make([]Element[T], len(init))
		for i, val := range init {
			elems[i] = Element[T]{
				Value: val,
			}
			if i < len(elems)-1 {
				elems[i].next = &elems[i+1]
			}
		}
		s.current = &elems[0]
	}
	return s
}

// Stack is a generic stack datastructure
type Stack[T any] struct {
	current *Element[T]
}

// Pop removes and returns the top Element from the stack
func (l *Stack[T]) Pop() *Element[T] {
	this := l.current
	if this != nil {
		l.current = this.next
	}
	return this
}

// Peek returns the top Element from the stack without removing it
func (l *Stack[T]) Peek() *Element[T] {
	return l.current
}

// Push adds a new item to the top of the stack
func (l *Stack[T]) Push(v T) {
	l.current = &Element[T]{next: l.current, Value: v}
}

// Swap replaces the current items value
func (l *Stack[T]) Swap(v T) {
	l.current.Value = v
}
