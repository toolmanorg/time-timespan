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

import "strconv"

type coefficient []rune

func newCoefficient() coefficient {
	return coefficient(make([]rune, 0))
}

// String is to implement the fmt.Stringer interface
func (c *coefficient) String() string {
	return string(*c)
}

func (c *coefficient) appendRune(r rune) (bool, *timespanErr) {
	if r == '-' || r == '+' {
		if len(*c) > 0 {
			return false, timespanError(misplacedSignErr, "misplaced '%c' in coefficient %q", r, c)
		}
		*c = append(*c, r)
		return true, nil
	}

	if r >= '0' && r <= '9' {
		*c = append(*c, r)
		return true, nil
	}

	return false, nil
}

func (c *coefficient) value(sign int) (int, *timespanErr) {
	if len(*c) < 1 {
		return 0, timespanError(missingCoefErr, "missing coefficient")
	}

	cs := string(*c)
	cv, err := strconv.Atoi(cs)
	if err != nil {
		return 0, timespanError(unparseableCoefErr, "unparseable coefficient: %q", cs)
	}

	// If a) the sign override is negative and
	//    b) the parsed value is positive and
	//    c) the string isn't explicitly positive,
	// negate the return value
	if sign < 0 && cv > 0 && cs[0] != '+' {
		cv *= -1
	}

	return cv, nil
}
