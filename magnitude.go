/*
Copyright 2017 Timothy E. Peoples

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
