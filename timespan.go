// Copyright 2017 Timothy E. Peoples. All rights reserved.
//
// TODO(tep): Add requisit licensing preamble

// Package timespan provides a broad-scaled extension to time.Duration.
//
// This package provides two types:
//
// 		Timespan - An extension to time.Duration
// 		Time     - an alias for time,Time (to ease use of Timespan values)
//
// Most of the functionality provided here is available using functions from
// the standard "time" package and is provided here merely as a convenience.
//
// However, this package's raison d'être is the function ParseTimespan, which
// provides the ability to specify a wide variety of broad-scaled time spans --
// from nanoseconds to many years -- as a simple, string value similar to that
// parseable by time.ParseDuration.
//
// For example, a time span of "1 year, 6 months" is specified as "1Y6M"
// or, its virtual equivalent, "18 months" as a simple "18M".  Timespan
// strings can be as simple as "3W" for "3 weeks" or something crazy like
// "1Y2M3W4D5h6m7s89ms" which is (hopefully) quite self explanatory.
//
// The full syntax for TimeSpan strings is described below.
//
package timespan

import (
	"fmt"
	"time"
)

// A Timespan represents a span of time with wide and varying resolutions.
//
// Unlike the standard time.Duration, which only provides accuracy and
// parseability at resolutions less than a day, a Timespan may cover many days,
// weeks, months or even years. It also encapsulates a time.Duration value to
// allow for resolutions as small a nanosecond.
//
// The abiguities around the length of a day or month (or even year) that
// restrict the broader scope of time.Duration are not a problem with Timespan
// since it stores each, distinctly varying magnitude separately.
//
// That said, a Timespan by itself is inherantly abiguous and only acquires
// precision when considered in the context of a specific point in time.
// Because of this, two separate Timespan values may be functionally equivalent
// from the perspective of one point in time but not from another.
//
// As an example, consider the Timespan values of "2 days" and "48 hours". In
// most cases, these two are functionally equivalent; 48 hours is always 48
// hours yet 2 days is usually 48 hours -- but not always.  In the approach to
// a daylight savings time cutover, the "2 days" value would be either 47 or 49
// hours (depending on whether we're "springing forward" or "falling back").
//
// In most cases however, these ambiguities are understood at the human level
// and Timespan will behave as the user intends without much further thought.
//
type Timespan struct {
	// Years in this Timespan
	Years int

	// Months in this Timespan
	Months int

	// Days in this Timespan
	Days int

	// A time.Duration for finer resolutions
	Duration time.Duration
}

// ParseTimespan returns a pointer to a new Timespan value which is the result
// of parsing a string representation for the desired Timespan.  The string is
// the conjunction of one or more coefficient+magnitude pairs plus an optional
// time.Duration string.
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
// case commutes across successive values until it is reversed. In other words,
// the stated positive or negative sense is always "sticky" for later values.
//
// Or rather, the positive or negative sense of each coefficient may be implied
// or expressly stated. By default, values are implied to be positive until one
// is explicitly declared to be negative. Subsequent values are then implied to
// be negative until an explicit positive coefficient is encountered.
//
// For example, "-1Y2M" is parsed to {Years: -1, Months: -2}. If your intent is
// instead for "Years" to be negative and "Months" to be positive (e.g. {Years:
// -1, Months: +2}), you must explicitly change the sense back to positive with
// "-1Y+2M".
//
// Zero value magnitudes may be omitted. However, each specified magnitide must
// be accompanied by a coefficient.  No whitespace is allowed anywhere in the
// string.  If supplied, the optional time.Duration string must be at the end
// of the string and must (of course) be parseable by time.ParseDuration.
//
// Since a week is always 7 days, the available "W" magnitude is provided
// merely as a convenience; it is not stored as part of the Timespan value.
// Coeffients provided in weeks are stored as multiples of 7 days.
//
// If ParseTimespan is unable to parse the given string, it returns nil and an
// approprate error.
//
// Finally, for those so inclined, the formal grammar for a Timespan string
// is shown in the following Pseudo-BNF:
//
//     <timespan>    := <periods>
//                    | <duration>
//                    | <periods> <duration>
// 
//     <periods>     := <period>
//                    | <periods> <period>
// 
//     <period>      := <coefficient> MAGNITUDE
//                    | SENSE <coefficient> MAGNITUDE
// 
//     <duration>    := {Anything parseable by time.ParseDuration}
// 
//     <coefficient> := DIGIT
//                    | <coefficient> DIGIT
// 
//     DIGIT         := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
// 
//     SENSE         := '-' | '+'
// 
//     MAGNITUDE     := 'Y' | 'M' | 'W' | 'D' | 'd'
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

// From returns the time.Time that results from applying the Timespan ts to the
// point in time t.  This is shorthand for:
//
// 		t.AddDate(ts.Years, ts.Months, ts.Days).Add(ts.Duration)
//
func (ts *Timespan) From(t time.Time) time.Time {
	return t.AddDate(ts.Years, ts.Months, ts.Days).Add(ts.Duration)
}

// Add returns a new *Timespan that is result of adding each member of ots to
// its corresponding member in ts. No combining, reduction or carry-over is
// performed.
//
// For example, if you add two Timespan values of 8 and 9 months, the result is
// always a Timespan value of 17 months (never 1 Year, 5 Months).
//
func (ts *Timespan) Add(ots *Timespan) *Timespan {
	return &Timespan{
		Years:    ts.Years + ots.Years,
		Months:   ts.Months + ots.Months,
		Days:     ts.Days + ots.Days,
		Duration: ts.Duration + ots.Duration,
	}
}

// Equal determines whether two Timespans are exactly equivalent to each other.
// Each member in ts is compared to its corresponding member in ots and all must
// be equivalent for Equal to return true.
//
// The Timespan values of "2 Days" and "48 Hours" are never equivalent in this
// context.
//
func (ts *Timespan) Equal(ots *Timespan) bool {
	return ts.Duration == ots.Duration &&
		ts.Days == ots.Days &&
		ts.Months == ots.Months &&
		ts.Years == ots.Years
}

// EqualAt determines whether two Timespans are functionally equivalent.  The
// Timespan values ts and ots are each evaluated at Time t and the result of
// each is compared. EqualAt returns true iff the two evluations resolve to the
// same point in time.
//
func (ts *Timespan) EqualAt(ots *Timespan, t time.Time) bool {
	return ts.From(t).Sub(t) == ots.From(t).Sub(t)
}

// Time is a convenience alias for time.Time provided simply to act as
// a receiver for the methods below.
//
type Time time.Time

// Add returns a new Time value after applying the given Timespan
//
func (t Time) Add(ts *Timespan) Time {
	return Time(ts.From(time.Time(t)))
}

// TimespansEqual compares the two Timespan values in the context of this Time.
// This is the same as:
//
// 		ts1.EqualAt(ts2, time.Time(t))
//
func (t Time) TimespansEqual(ts1, ts2 *Timespan) bool {
	return ts1.EqualAt(ts2, time.Time(t))
}

// String is shorthand for time.Time(t).String() and is provided to implement
// the fmt.Stringer interface.
//
func (t Time) String() string {
	return time.Time(t).String()
}
