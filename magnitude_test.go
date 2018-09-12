package timespan

import (
	"fmt"
	"strings"
	"testing"
)

type element struct {
	glyph rune
	magn  *magnitude
}

func be(g rune, l string) *element {
	return &element{glyph: g, magn: &magnitude{label: l, isSet: false, value: 0}}
}

func ve(g rune, l string, v int) *element {
	return &element{glyph: g, magn: &magnitude{label: l, isSet: true, value: v}}
}

func mkMagset(mags ...*element) magset {
	ms := magset(map[rune]*magnitude{})
	for _, m := range mags {
		ms[m.glyph] = m.magn
	}
	return ms
}

// String implements the fmt.Stringer interface
// to present magsets in a human readable form in testing errors
func (ms magset) String() string {
	ml := []string{}
	for g, m := range ms {
		ml = append(ml, fmt.Sprintf("'%c': {label: %q, isSet: %v, value: %d}", g, m.label, m.isSet, m.value))
	}
	return fmt.Sprintf("magset{%s}", strings.Join(ml, ", "))
}

// A utility functions for comparing two magsets
// (here since it's used only for testing)
func (ms magset) equal(oms magset) bool {
	for g, m := range ms {
		om, ok := oms[g]
		if !ok {
			return false
		}

		if !(m.label == om.label && m.isSet == om.isSet && m.value == om.value) {
			return false
		}
	}

	return true
}

// A quick test to ensure the 'equal()' function above actually works properly
func TestMagsetEqual(t *testing.T) {
	m1 := mkMagset(be('A', "one"), ve('B', "two", 42))
	m2 := mkMagset(be('A', "one"), ve('B', "two", 42))

	if !m1.equal(m2) {
		t.Errorf("magset.equal() failed to recognize identical magsets:\n\t1: %v\n\t2: %v", m1, m2)
	}

	m2['A'].set(1)

	if m1.equal(m2) {
		t.Errorf("magset.equal() failed to recognize dissimilar magsets:\n\t1: %v\n\t2: %v", m1, m2)
	}
}

//----- Now, we can get down to testing in earnest --------------------

func TestNewMagset(t *testing.T) {
	want := mkMagset(be('Y', "year"), be('M', "month"), be('W', "week"), be('D', "day"))
	got := newMagset()

	if !got.equal(want) {
		t.Errorf("newMagset() generated an invalid magset:\n\t Got: %v\n\tWant: %v", got, want)
	}
}

func TestMagsetGood(t *testing.T) {
	got := newMagset()

	got.set('Y', 1)
	got.set('M', 2)
	got.set('W', 3)
	got.set('D', -4)

	want := mkMagset(ve('Y', "year", 1), ve('M', "month", 2), ve('W', "week", 3), ve('D', "day", -4))

	if !got.equal(want) {
		t.Errorf("magset mismatch:\n\t Got: %v\n\tWant: %v", got, want)
	}
}

func TestMagsetBadRune(t *testing.T) {
	ms := mkMagset(be('A', "one"), ve('B', "two", 42))

	if err := ms.set('C', 1); err.errorType != unrecognizedMagErr {
		t.Errorf("attempting to set an unknown magnitude: bad error type: Got:%v Wanted:%v", err.errorType, unrecognizedMagErr)
	}
}

func TestMagsetUnkownOrder(t *testing.T) {
	ms := mkMagset(be('A', "one"), ve('B', "two", 42))

	err := ms.set('A', 1)
	if err == nil {
		t.Errorf("attempting to set magnitude of unknown order: no error returned when expected: wanted %v", magnOrderUnkownErr)
	} else if err.errorType != magnOrderUnkownErr {
		t.Errorf("attempting to set magnitude of unknown order: bad error type: Got:%v  Wanted:%v", err.errorType, magnOrderUnkownErr)
	}
}

func TestMagset(t *testing.T) {
	for i1, r1 := range magOrder {
		for i2, r2 := range magOrder {
			ms := newMagset()

			if err := ms.set(r1, i1); err != nil {
				t.Errorf("setting initial magnitude value for %q: unexpected error: Got:%v Wanted:%v", ms[r1].label, err, nil)
				continue
			} else {
				if v := ms.get(r1); v != i1 {
					t.Errorf("bad value from magset.get(%q): Got:%d Wanted:%d", ms[r1].label, v, i1)
				}
			}

			switch {
			case i1 == i2:
				if err := ms.set(r2, 2); err.errorType != magnRestatedErr {
					t.Errorf("restating magnitude value for %q: bad error type: Got:%v Wanted:%v", ms[r2].label, err, magnRestatedErr)
				}

			case i1 < i2:
				if err := ms.set(r2, i2); err != nil {
					t.Errorf("setting magnitude value for %q after %q: unexpected error: Got:%v Wanted:%v", ms[r2].label, ms[r1].label, err, nil)
				}

			case i1 > i2:
				err := ms.set(r2, i2)
				if err == nil {
					t.Errorf("setting magnitude value for %q after %q: no error returned when expected: wanted %v", ms[r2].label, ms[r1].label, magnOutOfOrderError)
				} else if err.errorType != magnOutOfOrderError {
					t.Errorf("setting magnitude value for %q after %q: bad error type: Got:%v  Wanted:%v", ms[r2].label, ms[r1].label, err, magnOutOfOrderError)
				}
			}
		}
	}
}
