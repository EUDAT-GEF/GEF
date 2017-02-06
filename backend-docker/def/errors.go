package def

import (
	"bytes"
	"fmt"
)

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

func Err(cause error, format string, a ...interface{}) Error {
	return Error{
		fmt.Sprintf(format, a...),
		cause,
	}
}
