package parser

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet[int](1, 2, 3)

	if !s.Contains(1) {
		t.Errorf("Set should contain 1")
	}

	if s.Contains(4) {
		t.Errorf("Set should not contain 4")
	}

	s.Add(4)

	if !s.Contains(4) {
		t.Errorf("Set should contain 4")
	}

	s.Remove(4)

	if s.Contains(4) {
		t.Errorf("Set should not contain 4")
	}

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(s.Members(), expected) {
		t.Errorf("Set members should be %v, but got %v", expected, s.Members())
	}
}
