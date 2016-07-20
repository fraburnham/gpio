package gpio

import "fmt"

type RpiGPIOError struct {
	msg string
}

func (e *RpiGPIOError) Error() string {
	return e.msg
}

func attachErrorCause(msg string, cause error) error {
	return &RpiGPIOError{msg: fmt.Sprintf("%s: %s", msg, cause.Error())}
}
