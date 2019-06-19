package mytmplate

import "errors"

// template type Stack(A)
type A int

type Stack struct {
	data []*A
}

func (s *Stack) Push(value *A) {
	s.data = append(s.data, value)
}

func (s *Stack) Pop() (*A, error) {
	length := len(s.data)
	if length == 0 {
		return nil, errors.New("Stack is empty")
	}
	value := s.data[length-1]
	s.data = s.data[:length-1]
	return value, nil
}

func NewStack() *Stack {
	return &Stack{
		data: []*A{},
	}
}
