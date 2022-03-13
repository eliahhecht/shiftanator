package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
)

type shift struct {
	a, b, c string
}

type schedule struct {
	shifts map[int]shift
}

func newSchedule() schedule {
	var s schedule
	s.shifts = make(map[int]shift)
	return s
}

func (sc schedule) add_shift(start int, sh shift) (schedule, error) {
	updated := sc
	fmt.Printf("Attempting to add %v:%v to %v\n", start, sh, sc)
	_, hasKey := updated.shifts[start]
	if hasKey {
		return schedule{}, errors.New(fmt.Sprintf("Shift %v already exists on this schedule", start))
	}
	updated.shifts[start] = sh
	return updated, nil
}

func (sc schedule) checkBalance(shiftAccessor func(shift) string) bool {
	var selectedShifts = make([]string, 6)
	for _, v := range sc.shifts {
		selectedShifts = append(selectedShifts, shiftAccessor(v))
	}
	sort.Strings(selectedShifts)
	return reflect.DeepEqual(
		selectedShifts, []string{"e", "e", "k", "k", "l", "l"})
}

func (sc schedule) isBalanced() bool {
	// fmt.Printf("checking balance for %v\n", sc)
	var accessors = []func(shift) string{
		func(s shift) string { return s.a },
		func(s shift) string { return s.b },
		func(s shift) string { return s.c },
	}
	for _, a := range accessors {
		if !sc.checkBalance(a) {
			return false
		}
	}
	return true
}

func (sc schedule) hasDoubleBreaks() bool {
	keys := make([]int, 0, len(sc.shifts))
	for k := range sc.shifts {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	lastBreak := "x"
	for _, k := range keys {
		if sc.shifts[k].c == lastBreak {
			return true
		}
		lastBreak = sc.shifts[k].c
	}
	return false
}

var shift_perms = []shift{
	{"e", "k", "l"},
	{"e", "l", "k"},
	{"l", "e", "k"},
	{"l", "k", "e"},
	{"k", "e", "l"},
	{"k", "l", "e"},
}

func schedule_perms() ([]schedule, error) {
	var start_times = []int{8, 10, 12, 14, 16, 18}
	var schedules = []schedule{newSchedule()}

	for _, start := range start_times {
		var updated = make([]schedule, 0)
		for _, sched := range schedules {
			for _, shift := range shift_perms {
				fmt.Printf("Evaluating shift perm %v\n", shift)
				updated_sched, err := sched.add_shift(start, shift)
				if err != nil {
					return nil, err
				}
				updated = append(updated, updated_sched)
			}
		}
		schedules = updated
	}

	return schedules, nil
}

func main() {
	perms, err := schedule_perms()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rand.Shuffle(len(perms), func(i, j int) {
		perms[i], perms[j] = perms[j], perms[i]
	})

	valid := make([]schedule, 0)

	for _, sc := range perms {
		if sc.isBalanced() && !sc.hasDoubleBreaks() {
			valid = append(valid, sc)
		}
	}

	fmt.Println(valid)
}
