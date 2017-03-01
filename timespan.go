// Copyright 2017 Timothy E. Peoples. All rights reserved.
//
// TODO(tep): Add requisit licensing preamble

// Package timespan provides a broad-scaled extension to time.Duration.
//
package timespan

// TODO(tep): Add more documentation!! XXX

import (
	"fmt"
	"time"
)

// A Timespan represents a span of time with wide and varying resolutions.
//
// Unlike the standard time.Duration, which only provides accuracy and
// parseability at resolutions less than a day, a Timespan may cover many days,
// weeks, months or even years. Additionally, it encapsulates a time.Duration
// value to also allow for resolutions as small a nanosecond.
//
// The abiguities around the length of a day or month (or even year) that
// restrict the broader scope of time.Duration are not a problem with Timespan
// since it stores several distinct magnitudes internally.
//
// That said, a Timespan by itself is inherantly abiguous and only acquires
// precision when considered in the context of a specific point in time.
// Because of this, two separate Timespan values may be functionally equivalent
// from the perspective of one point in time but not another.
//
// As an example, consider the Timespan values of "2 days" and "48 hours". In
// most cases, these two are functionally equivalent; 48 hours is always 48
// hours and 2 days is usually 48 hours.  However, in the approach to a DST
// cutover, the "2 days" value would be either 47 or 49 hours (depending on
// whether we're "springing forward" or "falling back").
//
// In most cases however, these ambiguities are understood at the human level
// and Timespan will behave as the user intends without much further thought.
//
type Timespan struct {
	// The number of years in this Timespan
	Years int

	// The number of months in this Timespan
	Months int

	// The number of days in this Timespan
	Days int

	// A fine-grained, sub-day duration
	Duration time.Duration
}

// ParseTimespan creates a new Timespan value by parsing a string representation
// for the desired Timespan.  This string is the conjunction of one or more
// coefficient+magnitude pairs plus an optional string parseable by time.Duration.
//
// Each coefficient is a signed integer value while its magnitude is one of the
// following single-character indicators:
//
//  Y: Years
//  M: Months
//  W: Weeks
//  D: Days (or optionally 'd')
//
// Each successive coefficient+magnitude pair must be specified in a decreasing
// magnitude order. In other words, years must be specified before months and
// months must be specified before weeks, etc... If this ordering is not kept,
// ParseTimespan will emit an error.  Also, each magnitude must be distinct. If
// any magnitude resolution is restated (e.g. "3W+1W"), parsing will fail.
//
// The natural tendency while parsing a Timespan string is to assume a nagative
// case commutes across successive values until it is reversed.
//
// For example, "-1Y2M" is parsed as "-1 year -2 months". If your intent is
// instead "-1 year +2 months", you must explicitly change the case back to
// positive using "-1Y+2M".
//
// Zero value magnitudes may be omitted. However, each specified magnitide must
// be accompanied by a coefficient.  No whitespace is allowed anywhere in the
// string.  If supplied, the optional time.Duration string must be at the end
// of the string and must (of course) be parseable by time.ParseDuration.
//
// The "weeks" magnitude is provided as a convenience and is not stored as part
// of the Timespan value. Coeffients provided in weeks are stored as multiples
// of 7 days.
//
// If ParseTimespan is unable to parse the given string, it returns nil and an
// approprate error.
//
func ParseTimespan(s string) (*Timespan, error) {
	ts := &Timespan{}
	ms := newMagset()

	sign := 1
	valid := false
	coef := newCoefficient()

	for i, r := range s {
		if d, err := time.ParseDuration(s[i:]); err == nil {
			ts.Duration = d
			valid = true
			break
		}

		if ok, err := coef.appendRune(r); err != nil {
			return nil, err.withTimespan(s)
		} else if ok {
			continue
		}

		v, err := coef.value(sign)
		if err != nil {
			return nil, err.withTimespan(s)
		}

		if err := ms.set(r, v); err != nil {
			return nil, err.withTimespan(s)
		}

		if v < 0 {
			sign = -1
		} else {
			sign = 1
		}

		valid = true
		coef = newCoefficient()
	}

	if !valid {
		return nil, fmt.Errorf("no value derived for Timespan %q", s)
	}

	ts.Years = ms.get('Y')
	ts.Months = ms.get('M')
	ts.Days = ms.get('D')
	ts.Days += ms.get('W') * 7

	return ts, nil
}

// String renders a Timespan into a form parseable by ParseTimespan.
func (ts *Timespan) String() string {
	s := ""

	if ts.Years != 0 {
		s = fmt.Sprintf("%dY", ts.Years)
	}

	if ts.Months != 0 {
		s = fmt.Sprintf("%s%dM", s, ts.Months)
	}

	if ts.Days != 0 {
		s = fmt.Sprintf("%s%dD", s, ts.Days)
	}

	if ts.Duration != 0 {
		s = fmt.Sprintf("%s%v", s, ts.Duration)
	}

	return s
}

// From returns the result of applying a Timespan to a given point in time.
// This is shorthand for:
// 		t.AddDate(ts.Years, ts.Months, ts.Days).Add(ts.Duration)
//
func (ts *Timespan) From(t time.Time) time.Time {
	return t.AddDate(ts.Years, ts.Months, ts.Days).Add(ts.Duration)
}

// Add combines one Timespan with another rendering a third Timespan Each
// distinct element is added together but no combination or carry-over is
// performed.
//
// For example, if you add two Timespan values of 8 and 9 months, the result is
// a Timespan value of 17 months (not 1 Year, 5 Months).
func (ts *Timespan) Add(ots *Timespan) *Timespan {
	return &Timespan{
		Years:    ts.Years + ots.Years,
		Months:   ts.Months + ots.Months,
		Days:     ts.Days + ots.Days,
		Duration: ts.Duration + ots.Duration,
	}
}

// Equal determines whether two Timespans are exactly equivalent Each distinct
// element is compared and all must be equivalent for Equal to return true.
// The Timespan values of "2 Days" and "48 Hours" are never equivalent in this
// context.
func (ts *Timespan) Equal(ots *Timespan) bool {
	return ts.Duration == ots.Duration &&
		ts.Days == ots.Days &&
		ts.Months == ots.Months &&
		ts.Years == ots.Years
}

// EqualAt determines whether two Timespans are functionally equivalent.
// The Timespan values ts and ots are each evaluated at ctime and the result
// of each is compared. EqualAt returns true iff the two evluations resolve
// in the same point in time.
func (ts *Timespan) EqualAt(ots *Timespan, ctime time.Time) bool {
	return ts.From(ctime).Sub(ctime) == ots.From(ctime).Sub(ctime)
}

// Time is a convenience alias for time.Time provided simply to act as
// a receiver for the methods below.
type Time time.Time

// Add returns a new Time value after applying the given Timespan
func (t Time) Add(ts *Timespan) Time {
	return Time(ts.From(time.Time(t)))
}

// TimespansEqual compares the two Timespan values in the context of this Time.
func (t Time) TimespansEqual(ts1, ts2 *Timespan) bool {
	return ts1.EqualAt(ts2, time.Time(t))
}

// String is shorthand for time.Time(t).String() and is provided to implement
// the fmt.Stringer interface.
func (t Time) String() string {
	return time.Time(t).String()
}
