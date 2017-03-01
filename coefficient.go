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
	if r == '-' {
		if len(*c) > 0 {
			return false, timespanError(missplacedDashErr, "misplaced '-' in coefficient %q", c)
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

func (c *coefficient) value() (int, *timespanErr) {
	if len(*c) < 1 {
		return 0, timespanError(missingCoefErr, "missing coefficient")
	}

	cs := string(*c)
	cv, err := strconv.Atoi(cs)
	if err != nil {
		return 0, timespanError(unparseableCoefErr, "unparseable coefficient: %q", cs)
	}
	return cv, nil
}
