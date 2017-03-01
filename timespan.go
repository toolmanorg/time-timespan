package timespan

// TODO(tep): Add documentation!! XXX

import (
	"fmt"
	"time"
)

type Timespan struct {
	years  int
	months int
	days   int
	dur    time.Duration
}

func NewTimespan(years, months, days int, duration time.Duration) *Timespan {
	return &Timespan{
		years:  years,
		months: months,
		days:   days,
		dur:    duration,
	}
}

func ParseTimespan(s string) (*Timespan, error) {
	ts := &Timespan{}
	ms := newMagset()

	valid := false
	coef := newCoefficient()

	for i, r := range s {
		if d, err := time.ParseDuration(s[i:]); err == nil {
			ts.dur = d
			valid = true
			break
		}

		if ok, err := coef.appendRune(r); err != nil {
			return nil, err.withTimespan(s)
		} else if ok {
			continue
		}

		v, err := coef.value()
		if err != nil {
			return nil, err.withTimespan(s)
		}

		if err := ms.set(r, v); err != nil {
			return nil, err.withTimespan(s)
		}

		valid = true
		coef = newCoefficient()
	}

	if !valid {
		return nil, fmt.Errorf("no value derived for Timespan %q", s)
	}

	ts.years = ms.get('Y')
	ts.months = ms.get('M')
	ts.days = ms.get('D')
	ts.days += ms.get('W') * 7

	return ts, nil
}

func (ts *Timespan) String() string {
	s := ""

	if ts.years != 0 {
		s = fmt.Sprintf("%dY", ts.years)
	}

	if ts.months != 0 {
		s = fmt.Sprintf("%s%dM", s, ts.months)
	}

	if ts.days != 0 {
		s = fmt.Sprintf("%s%dD", s, ts.days)
	}

	if ts.dur != 0 {
		s = fmt.Sprintf("%s%v", s, ts.dur)
	}

	return s
}

func (ts *Timespan) Delta(t time.Time) time.Time {
	return t.AddDate(ts.years, ts.months, ts.days).Add(ts.dur)
}

func (ts *Timespan) Add(ots *Timespan) *Timespan {
	return &Timespan{
		years:  ts.years + ots.years,
		months: ts.months + ots.months,
		days:   ts.days + ots.days,
		dur:    ts.dur + ots.dur,
	}
}

func (ts *Timespan) Equal(ots *Timespan) bool {
	return ts.dur == ots.dur &&
		ts.days == ots.days &&
		ts.months == ots.months &&
		ts.years == ots.years
}
