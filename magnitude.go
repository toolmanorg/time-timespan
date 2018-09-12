package timespan

import "strings"

const magOrder = "YMWD"

type magnitude struct {
	// glyph rune
	label string
	isSet bool
	value int
}

func (m *magnitude) set(val int) {
	m.value = val
	m.isSet = true
}

type magset map[rune]*magnitude

func newMagset() magset {
	return magset(map[rune]*magnitude{
		'Y': {label: "year"},
		'M': {label: "month"},
		'W': {label: "week"},
		'D': {label: "day"},
	})
}

func (ms magset) get(r rune) int {
	return ms[r].value
}

func (ms magset) set(r rune, val int) *timespanErr {
	if r == 'd' {
		r = 'D'
	}

	m, ok := ms[r]
	if !ok {
		return timespanError(unrecognizedMagErr, "unrecognized magnitude: %q", string(r))
	}

	i := strings.Index(magOrder, string(r))
	if i < 0 {
		return timespanError(magnOrderUnkownErr, "indeterminate order for magnitude: %q", string(r))
	}

	if m.isSet {
		return timespanError(magnRestatedErr, "magnitude %c restated (current:%d%c previous:%d%c)", r, val, r, m.value, r)
	}

	for _, g := range magOrder[i:] {
		om := ms[g]
		if om.isSet {
			return timespanError(magnOutOfOrderError, "magnitude out of order: %s specified before %s", om.label, m.label)
		}
	}

	m.set(val)

	return nil
}
