package gpio

import "time"

type InterruptEvent struct {
	Value     int
	Timestamp time.Time
}

type GPIO interface {
	MakeOutput() error
	MakeInput() error
	WriteValue(int) error
	ReadValue() (int, error)
	SetInterrupt(string, chan InterruptEvent) (chan bool, error)
	ClearInterrupt(chan bool) error
	Close() error
}
