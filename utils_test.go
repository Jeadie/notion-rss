package main

import (
	"fmt"
	"testing"
)

// ShouldPanic
func ShouldPanic(t *testing.T, f func(), expectedErrorMessage string, errMsg string) {
	t.Helper()
	defer func() {
		err := recover()
		if err.(error).Error() != expectedErrorMessage {
			t.Errorf(errMsg)
		}
	}()
	f()
	t.Errorf(errMsg)
}

// ShouldNotPanic
func ShouldNotPanic(t *testing.T, f func(), errMsg string) {
	t.Helper()
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf(errMsg)
		}
	}()
	f()
}

func TestPanicOnErrors(t *testing.T) {
	type TestCase struct {
		Errors          []error
		IsPanicExpected bool
		PanicErrMsg     string
		FailMsg         string
	}
	tests := []TestCase{
		{
			Errors:          []error{},
			IsPanicExpected: false,
			FailMsg:         "should not panic when no errors provided",
		},
		{
			Errors:          []error{fmt.Errorf("single Error")},
			IsPanicExpected: true,
			PanicErrMsg:     "single Error",
			FailMsg:         "For a single error, should panic with specific message",
		},
		{
			Errors:          []error{nil, fmt.Errorf("single Error")},
			IsPanicExpected: true,
			PanicErrMsg:     "single Error",
			FailMsg:         "should panic if input contains nil errors",
		},
		{
			Errors:          []error{fmt.Errorf("single Error"), fmt.Errorf("one more")},
			IsPanicExpected: true,
			PanicErrMsg:     "Multiple errors occured. Check output for details",
			FailMsg:         "For multiple errors, should provide generic message",
		},
	}

	fmt.Println("The following output may contain error messages. These are expected")
	for _, test := range tests {
		if test.IsPanicExpected {
			ShouldPanic(t, func() { PanicOnErrors(test.Errors) }, test.PanicErrMsg, test.FailMsg)
		} else {
			ShouldNotPanic(t, func() { PanicOnErrors(test.Errors) }, test.FailMsg)
		}
	}
	// Separate out expected error messages from rest of testing output
	fmt.Printf("End of expected errors\n\n")
}
