package gpio

import "fmt"

type GPIOError struct {
	msg string
}

func (e *GPIOError) Error() string {
	return e.msg
}

func NewError(msg string) *GPIOError {
	return &GPIOError{msg: msg}
}

func AttachErrorCause(msg string, cause error) error {
	return &GPIOError{msg: fmt.Sprintf("%s: %s", msg, cause.Error())}
}
