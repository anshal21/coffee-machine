package expressions

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

func (s *stack) Top() interface{} {
	if s.index == -1 {
		return nil
	}
	return s.elements[s.index]
}

func (s *stack) Pop() {
	if s.index >= 0 {
		s.index--
	}
}

func (s *stack) Push(val interface{}) {
	s.index++
	if s.index < len(s.elements) {
		s.elements[s.index] = val
		return
	} else {
		s.elements = append(s.elements, val)
	}
}
