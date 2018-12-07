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
	badDurationErr
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
