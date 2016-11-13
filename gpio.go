// The gpio package provides an interface for interacting with GPIO pins. (rpi, bbb, etc)
package gpio

import "time"

type InterruptEvent struct {
	Value     int
	Timestamp time.Time
	Err       error
}

type GPIO interface {
	MakeOutput() error
	MakeInput() error
	WriteValue(int) error
	ReadValue() (int, error)
	SetInterrupt(string, chan InterruptEvent, int) error
	ClearInterrupt() error
	Close() error
}
