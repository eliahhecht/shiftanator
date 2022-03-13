package main

import "testing"

func TestAddShiftDoesNotAffectOriginal(t *testing.T) {
	sched := newSchedule()
	sched.addShift(8, shift{})
	_, err := sched.addShift(8, shift{})
	if err != nil {
		t.Errorf("Expected no error but got %q", err)
	}
}

func TestHasDoubleBreaks_DoesntHaveThem(t *testing.T) {
	sched, err := newSchedule().addShift(8, shift{"k", "e", "l"})
	if err != nil {
		t.Error(err)
	}
	sched, err = sched.addShift(10, shift{"l", "e", "k"})
	if err != nil {
		t.Error(err)
	}

	double := sched.hasDoubleBreaks()
	if double {
		t.Errorf("Expected hasDoubleBreaks to return false for %v, but it returned true", sched)
	}
}

func TestBalanced_IsBalanced(t *testing.T) {
	starts := []int{8, 10, 12, 14, 16, 18}
	shifts := []shift{
		{"k", "e", "l"},
		{"e", "l", "k"},
		{"l", "k", "e"},
		{"k", "e", "l"},
		{"e", "l", "k"},
		{"l", "k", "e"},
	}
	sched := newSchedule()
	for i, s := range starts {
		var err error
		sched, err = sched.addShift(s, shifts[i])
		if err != nil {
			t.Error(err)
		}
	}

	balanced := sched.isBalanced()
	if !balanced {
		t.Errorf("Expected isBalanced to return true for %v, but it returned false", sched)
	}

}
