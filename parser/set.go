package parser

import (
	"strings"
	"unicode/utf8"
)

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

// maxFoldLen bounds the stack buffer used by lookupFold and containsFold;
// longer inputs fall back to strings.ToUpper.
const maxFoldLen = 64

// upperASCII writes the upper-cased form of s into buf and reports whether it
// could do so without allocating. It refuses non-ASCII input so callers can
// fall back to strings.ToUpper and keep identical semantics.
func upperASCII(s string, buf *[maxFoldLen]byte) bool {
	if len(s) > len(buf) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			return false
		}
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		buf[i] = c
	}
	return true
}

// lookupFold indexes m by the upper-cased form of s. Indexing with
// string(buf[:n]) lets the compiler skip the string allocation that
// strings.ToUpper would need for every lower-case identifier.
func lookupFold[V any](m map[string]V, s string) (V, bool) {
	var buf [maxFoldLen]byte
	if !upperASCII(s, &buf) {
		v, ok := m[strings.ToUpper(s)]
		return v, ok
	}
	v, ok := m[string(buf[:len(s)])]
	return v, ok
}

// containsFold reports whether the upper-cased form of s is a member of set.
func containsFold(set *Set[string], s string) bool {
	_, ok := lookupFold(set.m, s)
	return ok
}
