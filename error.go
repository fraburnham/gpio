package gpio

import "fmt"

type GPIOError struct {
	msg string
}

func (e *GPIOError) Error() string {
	return e.msg
}

// NewError creates a fresh GPIOError.
func NewError(msg string) *GPIOError {
	return &GPIOError{msg: msg}
}

// AttachErrorCause attaches details to an error caused by upstream code and returns a
// GPIOError.
func AttachErrorCause(msg string, cause error) error {
	return &GPIOError{msg: fmt.Sprintf("%s: %s", msg, cause.Error())}
}
