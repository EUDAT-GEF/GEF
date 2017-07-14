package def

import (
	"bytes"
	"fmt"
)

// Error is the main error type
type Error struct {
	message string
	cause   error
}

func (e Error) Error() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s", e.message)
	if e.cause != nil {
		fmt.Fprintf(&b, "\n\tcaused by: %s", e.cause.Error())
	}
	return b.String()
}

// Err fn creates a new Error from an optional existing cause
func Err(cause error, format string, a ...interface{}) Error {
	return Error{
		fmt.Sprintf(format, a...),
		cause,
	}
}

// PermissionError is returned when the currently loggedin user
// is denied an action
type PermissionError struct {
	message string
}

func (e PermissionError) Error() string {
	return e.message
}

// PermissionErr creates a new PermissionError
func PermissionErr(format string, a ...interface{}) PermissionError {
	return PermissionError{
		fmt.Sprintf(format, a...),
	}
}
