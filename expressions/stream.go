package expressions

type stream struct {
	s   []rune
	pos int
}

const (
	_EndOfStream = '\n'
)

func (s *stream) GetNext() rune {
	if s.pos >= len(s.s) {
		return _EndOfStream
	}
	s.pos++
	return s.s[s.pos-1]
}

func (s *stream) Position() int {
	return s.pos
}

func (s *stream) Rewind() {
	s.pos--
	if s.pos < 0 {
		s.pos = 0
	}
}

func newStream(s string) *stream {
	return &stream{
		s:   []rune(s),
		pos: 0,
	}
}
