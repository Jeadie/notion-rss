package main

import "fmt"

// PanicOnErrors prints all non-nil err in errors and panics if there is at least one non-nil
// error in errors. Otherwise, return normally.
func PanicOnErrors(errors []error) {
	// Only used if one error (for better error handling).
	var firstErr error
	errN := 0

	// Print all non-nil errors.
	for _, err := range errors {
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			errN++
			firstErr = err
		}
	}

	// Multiple errors, panic with generic message.
	if errN > 1 {
		panic(fmt.Errorf("Multiple errors occured. Check output for details"))
	}

	if errN == 1 {
		panic(firstErr)
	}
}
