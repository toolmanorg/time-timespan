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
	_, err := coef.value()
	if err == nil {
		t.Errorf("No error while retreiving value from empty coefficient: Wanted %s", missingCoefErr)
	} else if err != nil && err.errorType != missingCoefErr {
		t.Errorf("Incorrect error retreiving value from empty coefficient: Got %s; Wanted %s", err, missingCoefErr)
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
	if got, err := coef.value(); err != nil {
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

func TestCoefficientBadDash(t *testing.T) {
	coef := coefficient("123")

	// Append '-' character to populated coefficient
	// Expect: false, missplacedDashErr
	if ok, err := coef.appendRune('-'); err != nil {
		if err.errorType != missplacedDashErr {
			t.Errorf("Incorrect error returned while appending '-' to non-empty coefficient: Got %s; Wanted %s", err, missplacedDashErr)
		}
	} else if ok {
		t.Error("appendRune erroneously accepted '-' as a valid digit")
	}
}

func TestCoefficientGoodDash(t *testing.T) {
	coef := newCoefficient()

	if ok, err := coef.appendRune('-'); err != nil {
		t.Errorf("Error appending '-' to empty coefficient: %v", err)
	} else if !ok {
		t.Errorf("appendRune failed to accept '-' for empty coefficient")
	}

	coef = append(coef, '1', '2')

	if v, err := coef.value(); err != nil {
		t.Errorf("Error acquiring value of coefficient %q: %v", coef, err)
	} else if v != -12 {
		t.Errorf("Coefficient value mismatch: Got:%d Wanted:%d", v, -12)
	}
}
