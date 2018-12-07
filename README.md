
[![GoDoc](https://godoc.org/toolman.org/time/timespan?status.svg)](https://godoc.org/toolman.org/time/timespan) [![Go Report Card](https://goreportcard.com/badge/toolman.org/time/timespan)](https://goreportcard.com/report/toolman.org/time/timespan)

# timespan
`import "toolman.org/time/timespan"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="install">Install</a>

```bash

    go get toolman.org/time/timespan

```

## <a name="pkg-overview">Overview</a>
Package timespan provides a broad-scaled extension to time.Duration.

This package provides two types:

```text

	Timespan - An extension to time.Duration
	Time     - an alias for time.Time (to ease use of Timespan values)

```

Most of the functionality provided here is available using functions from
the standard "time" package and is provided here merely as a convenience.

However, this package's raison d'être is the function ParseTimespan, which
provides the ability to specify a wide variety of broad-scaled time periods
-- from nanoseconds to many years -- as a simple, string value similar to
that parseable by time.ParseDuration.

For example, a time span of "1 year, 6 months" is specified as "1Y6M"
or, its virtual equivalent, "18 months" as a simple "18M".  Timespan
strings can be as simple as "3W" for "3 weeks" or something crazy like
"1Y2M3W4D5h6m7s89ms" which is (hopefully) quite self explanatory.

### Motivation
Unlike the standard time.Duration, which only provides accuracy and
parseability at resolutions less than a day, a Timespan may cover many days,
weeks, months or even years. It also encapsulates a time.Duration value to
allow for resolutions as small a nanosecond.

The abiguities around the length of a day or month (or even year) that
restrict the broader scope of time.Duration are not a problem with Timespan
since it stores each, distinctly varying magnitude separately.

That said, a Timespan by itself is inherantly abiguous and only acquires
precision when considered in the context of a specific point in time.
Because of this, two separate Timespan values may be functionally equivalent
from the perspective of one point in time but not from another.

As an example, consider the Timespan values of "2 days" and "48 hours". In
most cases, these two are functionally equivalent; 48 hours is always 48
hours yet 2 days is sometimes not 48 hours.  In the approach to a daylight
savings time cutover, the "2 days" value would be either 47 or 49 hours
(depending on whether we're "springing forward" or "falling back").

In most cases however, these ambiguities are understood at the human level
and Timespan will behave as the user intends without much further thought.

### Parsing
A Timespan string is the conjunction of one or more periods (as
coefficient+magnitude pairs) plus an optional time.Duration string.

Parsing is governed by the following rules:

1. Each period is a coefficient magnitude pair where the coefficient is
a signed integer value and its magnitude is one of the following
single-character indicators:

```text

	Y: Years
	M: Months
	W: Weeks
	D: Days (or optionally 'd')

```

2. Each successive period must be specified in a decreasing magnitude order.
In other words, years must be specified before months and months before
weeks, etc... Periods specified out of order will cause an error.

3. The magnitude of each period must be distinct; any restated magnitude
(e.g.  "3W-1W"), causes a parsing error.

4. The positive or negative sense of each coefficient may be implied or
expressly stated. By default, values are assumed positive until one is
explicitly declared to be negative. Subsequent (implicitly signed) values
are assumed to be negative until an explicit positive coefficient is
encountered.

5. Zero value magnitudes may be omitted.

6. Each specified magnitide must be accompanied by a coefficient.

7. No whitespace is allowed anywhere in the string.

8. If supplied, the optional time.Duration string must be at the end of the
string and must (of course) be parseable by time.ParseDuration.

The natural tendency while parsing a Timespan string is to assume a nagative
sign commutes across successive values until it is reversed. In other words,
the stated positive or negative sign is always "sticky" for later values.

For example, "-1Y2M" is parsed to {Years: -1, Months: -2}. If your intent is
instead for "Years" to be negative and "Months" to be positive (e.g. {Years:
-1, Months: +2}), you must explicitly change the sign back to positive with
"-1Y+2M".

Since a week is always 7 days, the available "W" magnitude is provided
merely as a convenience; it is not stored as part of the Timespan value.
Coeffients provided in weeks are stored as multiples of 7 days.

If ParseTimespan is unable to parse the given string, it returns nil and an
approprate error.

### Grammar
Finally, for those so inclined, the formal grammar for a Timespan string
is shown in the following Pseudo-BNF:

```bnf

	<timespan>    := <periods>
	               | <duration>
	               | <periods> <duration>
	
	<periods>     := <period>
	               | <periods> <period>
	
	<period>      := <coefficient> MAGNITUDE
	               | SIGN <coefficient> MAGNITUDE
	
	<coefficient> := DIGIT
	               | <coefficient> DIGIT
	
	DIGIT         := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
	
	SIGN          := '-' | '+'
	
	MAGNITUDE     := 'Y' | 'M' | 'W' | 'D' | 'd'
	
	<duration>    := {Anything parseable by time.ParseDuration}

```


## <a name="pkg-index">Index</a>
* [type Time](#Time)
  * [func (t Time) Add(ts *Timespan) Time](#Time.Add)
  * [func (t Time) String() string](#Time.String)
  * [func (t Time) TimespansEqual(ts1, ts2 *Timespan) bool](#Time.TimespansEqual)
* [type Timespan](#Timespan)
  * [func ParseTimespan(s string) (*Timespan, error)](#ParseTimespan)
  * [func (ts *Timespan) Add(ots *Timespan) *Timespan](#Timespan.Add)
  * [func (ts *Timespan) Equal(ots *Timespan) bool](#Timespan.Equal)
  * [func (ts *Timespan) EqualAt(ots *Timespan, t time.Time) bool](#Timespan.EqualAt)
  * [func (ts *Timespan) From(t time.Time) time.Time](#Timespan.From)
  * [func (ts *Timespan) String() string](#Timespan.String)


#### <a name="pkg-files">Package files</a>
[coefficient.go](/src/toolman.org/timespan/coefficient.go) [errors.go](/src/toolman.org/timespan/errors.go) [errtype_string.go](/src/toolman.org/timespan/errtype_string.go) [magnitude.go](/src/toolman.org/timespan/magnitude.go) [timespan.go](/src/toolman.org/timespan/timespan.go) 






## <a name="Time">type</a> [Time](/src/target/timespan.go?s=9314:9333#L284)
``` go
type Time time.Time
```
Time is a convenience alias for time.Time provided simply to act as
a receiver for the methods below.










### <a name="Time.Add">func</a> (Time) [Add](/src/target/timespan.go?s=9404:9440#L288)
``` go
func (t Time) Add(ts *Timespan) Time
```
Add returns a new Time value after applying the given Timespan




### <a name="Time.String">func</a> (Time) [String](/src/target/timespan.go?s=9838:9867#L304)
``` go
func (t Time) String() string
```
String is shorthand for time.Time(t).String() and is provided to implement
the fmt.Stringer interface.




### <a name="Time.TimespansEqual">func</a> (Time) [TimespansEqual](/src/target/timespan.go?s=9628:9681#L297)
``` go
func (t Time) TimespansEqual(ts1, ts2 *Timespan) bool
```
TimespansEqual compares the two Timespan values in the context of this Time.
This is the same as:


	ts1.EqualAt(ts2, time.Time(t))




## <a name="Timespan">type</a> [Timespan](/src/target/timespan.go?s=5606:5813#L130)
``` go
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
```
A Timespan represents a span of time with wide and varying resolutions.







### <a name="ParseTimespan">func</a> [ParseTimespan](/src/target/timespan.go?s=5961:6008#L147)
``` go
func ParseTimespan(s string) (*Timespan, error)
```
ParseTimespan returns a pointer to a new Timespan value which is the result
of parsing a string representation for the desired Timespan.





### <a name="Timespan.Add">func</a> (\*Timespan) [Add](/src/target/timespan.go?s=8141:8189#L249)
``` go
func (ts *Timespan) Add(ots *Timespan) *Timespan
```
Add returns a new *Timespan that is result of adding each member of ots to
its corresponding member in ts. No combining, reduction or carry-over is
performed.

For example, if you add two Timespan values of 8 and 9 months, the result is
always a Timespan value of 17 months (never 1 Year, 5 Months).




### <a name="Timespan.Equal">func</a> (\*Timespan) [Equal](/src/target/timespan.go?s=8660:8705#L265)
``` go
func (ts *Timespan) Equal(ots *Timespan) bool
```
Equal determines whether two Timespans are exactly equivalent to each other.
Each member in ts is compared to its corresponding member in ots and all must
be equivalent for Equal to return true.

The Timespan values of "2 Days" and "48 Hours" are never equivalent in this
context.




### <a name="Timespan.EqualAt">func</a> (\*Timespan) [EqualAt](/src/target/timespan.go?s=9089:9149#L277)
``` go
func (ts *Timespan) EqualAt(ots *Timespan, t time.Time) bool
```
EqualAt determines whether two Timespans are functionally equivalent.  The
Timespan values ts and ots are each evaluated at Time t and the result of
each is compared. EqualAt returns true iff the two evluations resolve to the
same point in time.




### <a name="Timespan.From">func</a> (\*Timespan) [From](/src/target/timespan.go?s=7703:7750#L238)
``` go
func (ts *Timespan) From(t time.Time) time.Time
```
From returns the time.Time that results from applying the Timespan ts to the
point in time t.  This is shorthand for:


	t.AddDate(ts.Years, ts.Months, ts.Days).Add(ts.Duration)




### <a name="Timespan.String">func</a> (\*Timespan) [String](/src/target/timespan.go?s=7194:7229#L211)
``` go
func (ts *Timespan) String() string
```
String renders a Timespan into a form parseable by ParseTimespan.

