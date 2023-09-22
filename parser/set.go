package parser

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](members ...T) *Set[T] {
	m := make(map[T]struct{})
	for _, member := range members {
		m[member] = struct{}{}
	}
	return &Set[T]{m: m}
}

func (s *Set[T]) Add(member T) {
	s.m[member] = struct{}{}
}

func (s *Set[T]) Remove(member T) {
	delete(s.m, member)
}

func (s *Set[T]) Contains(member T) bool {
	_, ok := s.m[member]
	return ok
}

func (s *Set[T]) Members() []T {
	members := make([]T, 0, len(s.m))
	for member := range s.m {
		members = append(members, member)
	}
	return members
}
