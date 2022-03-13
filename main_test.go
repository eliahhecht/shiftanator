package main

import (
	"testing"
)

func TestAddShiftDoesNotAffectOriginal() {
	sched := newSchedule()
	updated = sched.add_shift(8, shift{})
	sched.add_shift(8, shift{})
}
