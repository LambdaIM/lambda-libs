package abi

import (
	"errors"
	"sync"
)

type Stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []interface{}
}

func NewStack() *Stack {
	return &Stack{sync.Mutex{}, make([]interface{}, 0)}
}

func (s *Stack) Push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *Stack) Pop() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]

	return res, nil
}
