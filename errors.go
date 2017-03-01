package timespan

//go:generate stringer -type=errType

import "fmt"

type errType int

const (
	noErr errType = iota
	misplacedSignErr
	missingCoefErr
	unparseableCoefErr
	unrecognizedMagErr
	magnOrderUnkownErr
	magnRestatedErr
	magnOutOfOrderError
)

type timespanErr struct {
	errorType errType
	tsString  string
	message   string
}

func timespanError(etype errType, mesg string, args ...interface{}) *timespanErr {
	return &timespanErr{
		errorType: etype,
		message:   fmt.Sprintf(mesg, args...),
	}
}

// Error is part of the error interface
func (te *timespanErr) Error() string {
	if te.tsString == "" {
		return te.message
	}

	return fmt.Sprintf("parsing Timespan %q: %s", te.tsString, te.message)
}

func (te *timespanErr) withTimespan(ts string) *timespanErr {
	te.tsString = ts
	return te
}
