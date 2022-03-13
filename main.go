package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
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

func (sc schedule) deepCopy() schedule {
	c := newSchedule()
	for k, v := range sc.shifts {
		c.shifts[k] = v
	}
	return c
}

func (sc schedule) addShift(start int, sh shift) (schedule, error) {
	// fmt.Printf("Attempting to add %v:%v to %v\n", start, sh, sc)
	newSched := sc.deepCopy()
	_, hasKey := sc.shifts[start]
	if hasKey {
		return schedule{}, fmt.Errorf("shift %v already exists on this schedule", start)
	}
	newSched.shifts[start] = sh
	return newSched, nil
}

func (sc schedule) checkBalance(shiftAccessor func(shift) string) bool {
	var selectedShifts = make([]string, 0)
	for _, v := range sc.shifts {
		selectedShifts = append(selectedShifts, shiftAccessor(v))
		// fmt.Printf("Selected shifts is %v\n", selectedShifts)
	}
	sort.Strings(selectedShifts)
	balanced := reflect.DeepEqual(
		selectedShifts, []string{"e", "e", "k", "k", "l", "l"})
	return balanced
}

func (sc schedule) String() string {
	sb := strings.Builder{}
	keys := sc.sortedKeys()
	for _, k := range keys {
		shift := sc.shifts[k]
		sb.WriteString(fmt.Sprintf("%2d %s%s%s\n", k,
			strings.ToUpper(shift.a), strings.ToUpper(shift.b), strings.ToUpper(shift.c)))
	}
	return sb.String()
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
			// fmt.Printf("Imbalanced sched: %v\n", sc)
			return false
		}
	}
	return true
}

func (sc schedule) sortedKeys() []int {
	keys := make([]int, 0, len(sc.shifts))
	for k := range sc.shifts {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (sc schedule) hasDoubleBreaks() bool {
	keys := sc.sortedKeys()

	lastBreak := "x"
	for _, k := range keys {
		if sc.shifts[k].c == lastBreak {
			// fmt.Printf("Found a double break in %q", sc)
			return true
		}
		lastBreak = sc.shifts[k].c
	}
	return false
}

func (sc schedule) hasDoubleAs() bool {
	keys := sc.sortedKeys()

	lastBreak := "x"
	for _, k := range keys {
		if sc.shifts[k].a == lastBreak {
			return true
		}
		lastBreak = sc.shifts[k].a
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
				// fmt.Printf("Evaluating shift perm %v\n", shift)
				updated_sched, err := sched.addShift(start, shift)
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

	// fmt.Println(perms)

	valid := make([]schedule, 0)

	lHasLastBreak := func(sc schedule) bool { return sc.shifts[18].c == "l" }
	firstShiftIsKle := func(sc schedule) bool { return sc.shifts[8] == shift{"k", "l", "e"} }
	eIsAvailableToCook := func(sc schedule) bool {
		return sc.shifts[16].a != "e" && sc.shifts[18].a != "e"
	}

	for _, sc := range perms {
		if sc.isBalanced() &&
			!sc.hasDoubleBreaks() &&
			lHasLastBreak(sc) &&
			firstShiftIsKle(sc) &&
			eIsAvailableToCook(sc) &&
			!sc.hasDoubleAs() {
			valid = append(valid, sc)
		}
	}

	for _, sc := range valid {
		fmt.Println(sc)
	}
	fmt.Println(len(valid))
}
