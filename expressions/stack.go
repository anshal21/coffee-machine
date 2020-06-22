package expressions

// Stack represents a stack datastructre and expose the methods
// related to it
type stack struct {
	elements []interface{}
	index    int
}

func newStack() *stack {
	return &stack{
		elements: make([]interface{}, 0),
		index:    -1,
	}
}

// Top gives the current stack top
func (s *stack) Top() interface{} {
	if s.index == -1 {
		return nil
	}
	return s.elements[s.index]
}

// Pop removes the current stack top
// does nothing if the stack is empty
func (s *stack) Pop() {
	if s.index >= 0 {
		s.index--
	}
}

// Push adds the element to the stack top
func (s *stack) Push(val interface{}) {
	s.index++
	if s.index < len(s.elements) {
		s.elements[s.index] = val
		return
	}

	s.elements = append(s.elements, val)
}
