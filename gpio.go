package gpio

import "time"

type InterruptEvent struct {
	value int
	timestamp time.Time
}

type GPIO interface {
	MakeOutput() error
	MakeInput() error
	WriteValue(int) error
	ReadValue() (int, error)
	Interrupt(string, chan InterruptEvent) error
	Close() error
}
