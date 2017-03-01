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
