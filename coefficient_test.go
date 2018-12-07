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

import (
	"strconv"
	"testing"
)

func TestCoefficientEmptyAppendNonDigit(t *testing.T) {
	coef := newCoefficient()

	// Append a non-digit to an empty coefficient
	// Expect: false, nil
	if ok, err := coef.appendRune('x'); err != nil {
		t.Errorf("Error appending 'x' to coefficient %q", coef)
	} else if ok {
		t.Error("appendRune erroneously accepted 'x' as a valid digit")
	}
}

func TestCoefficientEmptyValue(t *testing.T) {
	coef := newCoefficient()

	// Retrieve value from empty coefficient
	// Expect: _, missingCoefErr
	_, err := coef.value(1)
	if err == nil {
		t.Errorf("No error while retreiving value from empty coefficient: Wanted %v", missingCoefErr)
	} else if err != nil && err.errorType != missingCoefErr {
		t.Errorf("Incorrect error retreiving value from empty coefficient: Got %v; Wanted %v", err, missingCoefErr)
	}
}

func TestCoefficientValidDigits(t *testing.T) {
	want := 9876543210
	chars := strconv.Itoa(want)
	coef := newCoefficient()

	// Append each valid digit in turn
	// Expect: true, nil (for each digit)
	for _, r := range chars {
		if ok, err := coef.appendRune(r); err != nil {
			t.Errorf("Error appending '%c' to coefficient %q", r, coef)
		} else if !ok {
			t.Errorf("appendRune failed to accept valid character '%c'", r)
		}
	}
}

func TestCoefficientValue(t *testing.T) {
	want := 12345
	coef := coefficient(strconv.Itoa(want))

	// Retrieve value from populated coefficient
	// Expect: 12345, nil
	if got, err := coef.value(1); err != nil {
		t.Errorf("Error discerning value of coefficient %q: %v", coef, err)
	} else if got != want {
		t.Errorf("Value mismatch for coefficient %q: Got %d; Wanted %d", coef, got, want)
	}
}

func TestCoefficientNonDigit(t *testing.T) {
	coef := coefficient("123")
	ndc := 'x'

	// Append non-digit character to populated coefficient
	// Expect: false, nil
	if ok, err := coef.appendRune(ndc); err != nil {
		t.Errorf("Error appending '%c' to coefficient %q", ndc, coef)
	} else if ok {
		t.Errorf("coefficient.appendRune() erroneously accepted '%c' as a valid digit", ndc)
	}

}

func TestCoefficientBadSign(t *testing.T) {
	coef := coefficient("123")

	// Append '-' and '+' chars to populated coefficient
	// Expect: false, misplacedSignErr
	for _, c := range []rune{'-', '+'} {
		if ok, err := coef.appendRune(c); err != nil {
			if err.errorType != misplacedSignErr {
				t.Errorf("Incorrect error returned while appending '%c' to non-empty coefficient: Got %v; Wanted %v", c, err, misplacedSignErr)
			}
		} else if ok {
			t.Errorf("appendRune erroneously accepted '%c' as a valid digit", c)
		}
	}
}

func TestCoefficientGoodSign(t *testing.T) {
	coef := newCoefficient()

	if ok, err := coef.appendRune('-'); err != nil {
		t.Errorf("Error appending '-' to empty coefficient: %v", err)
	} else if !ok {
		t.Errorf("appendRune failed to accept '-' for empty coefficient")
	}

	coef = append(coef, '1', '2')

	if v, err := coef.value(1); err != nil {
		t.Errorf("Error acquiring value of coefficient %q: %v", coef, err)
	} else if v != -12 {
		t.Errorf("Coefficient value mismatch: Got:%d Wanted:%d", v, -12)
	}
}

func TestCoefficientSignOverride(t *testing.T) {
	for _, str := range []string{"123", "-123", "+123"} {
		for _, sign := range []int{1, -1} {
			coef := coefficient(str)
			if got, err := coef.value(sign); err != nil {
				t.Errorf("Error discerning value of coefficient %q: %v", coef, err)
			} else {
				var want int
				switch {
				case str[0] == '+':
					want = 123

				case str[0] == '-':
					want = -123

				default:
					want = 123 * sign
				}
				if got != want {
					t.Errorf("Coefficient sign override failure.  Got:%d Wanted:%d", got, want)
				}
			}
		}
	}
}
