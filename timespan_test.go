package timespan

import (
	"testing"
	"time"
)

type testdata struct {
	str   string
	want  *Timespan
	etype errType
}

func TestParseTimespanGood(t *testing.T) {
	data := []testdata{
		{str: "1h30m", want: &Timespan{0, 0, 0, 1*time.Hour + 30*time.Minute}},
		{str: "-1h30m", want: &Timespan{0, 0, 0, -1*time.Hour - 30*time.Minute}},
		{str: "2D1h", want: &Timespan{0, 0, 2, 1 * time.Hour}},
		{str: "4W1d", want: &Timespan{0, 0, 29, 0}},
		{str: "4W-1d", want: &Timespan{0, 0, 27, 0}},
		{str: "1Y2M3W4D5h6m7s", want: &Timespan{1, 2, 25, 5*time.Hour + 6*time.Minute + 7*time.Second}},
	}

	for _, td := range data {
		got, err := ParseTimespan(td.str)
		if err != nil {
			t.Error(err)
			continue
		}

		if !(td.want.years == got.years && td.want.months == got.months && td.want.days == got.days && td.want.dur == got.dur) {
			t.Errorf("Mismatch parsing Timespan %q  Got:%+v  Wanted:%+v", td.str, got, td.want)
		}
	}
}

func TestParseTimespanBad(t *testing.T) {
	data := []testdata{
		{str: "Y", etype: missingCoefErr},
		{str: "1h2D", etype: unrecognizedMagErr},
		{str: "1D2W", etype: magnOutOfOrderError},
		{str: "3W2W", etype: magnRestatedErr},
		{str: "4W1-D", etype: missplacedDashErr},
	}

	for _, td := range data {
		_, err := ParseTimespan(td.str)
		if err == nil {
			t.Errorf("No error found parsing invalid Timespan %q: wanted:%v", td.str, td.etype)
			continue
		}

		if tse, ok := err.(*timespanErr); !ok {
			t.Fatalf("Error returned while parsing invalid Timespan %q is not a timespanError!! (...I dunno what's going on :-/ )", td.str)
		} else {
			if tse.errorType != td.etype {
				t.Errorf("Error mismatch parsing invalid Timespan %q: got %v; wanted %v", td.str, tse.errorType, td.etype)
			}
		}
	}
}

func TestTimespanEqual(t *testing.T) {
	ts1 := &Timespan{1, 2, 3, 4 * time.Hour}
	ts2 := &Timespan{1, 2, 3, 4 * time.Hour}

	if !ts1.Equal(ts2) {
		t.Errorf("Timespans should be identical:\n\t 1) %v\n\t2) %v")
	}
}

func TestTimespanString(t *testing.T) {
	ts, err := ParseTimespan("2Y2M2W2h30m")
	if err != nil {
		t.Fatal(err)
	}

	want := "2Y2M14D2h30m0s"
	got := ts.String()

	if got != want {
		t.Errorf("Timespan rendered to improper string:\n\t Got: %q\n\tWant: %q", got, want)
	}
}

func TestTimespanDelta(t *testing.T) {
	ts := &Timespan{0, 2, 14, 2*time.Hour + 30*time.Minute}

	base := time.Date(2014, 03, 03, 17, 0, 0, 0, time.UTC)
	want := time.Date(2014, 05, 17, 19, 30, 0, 0, time.UTC)

	got := ts.Delta(base)

	if !want.Equal(got) {
		t.Errorf("Timespan Delta Mismatch:\n\t Got: %v\n\tWant: %v", got, want)
	}
}

func TestTimespanAdd(t *testing.T) {
	ts1 := &Timespan{0, 2, 14, 2*time.Hour + 30*time.Minute}
	ts2 := &Timespan{1, 1, 10, 3*time.Hour + 30*time.Minute}

	want := &Timespan{1, 3, 24, 6 * time.Hour}

	got := ts1.Add(ts2)

	if !want.Equal(got) {
		t.Errorf("Timespan Add mismatch:\n\t Got: %+v\n\tWant: %+v", got, want)
	}
}
