package gpio

type GPIO interface {
	MakeOutput() error
	MakeInput() error
	WriteValue(int) error
	ReadValue() (int, error)
	Close() error
}
