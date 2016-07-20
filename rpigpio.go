package gpio

import (
	"fmt"
	"os"
)

type RpiGPIO struct {
	pin      int
	isOutput bool
	isExported bool
	baseDir  string
}

func write(file string, data string) error {
	f, err := os.OpenFile(file, os.O_WRONLY, 0644)
	defer f.Close()

	if err == nil {
		_, err = f.WriteString(data)
	}

	return err
}

func exportPin(baseDir string, pin int) error {
	return write(fmt.Sprintf("%s/export", baseDir),
		fmt.Sprintf("%d", pin))
}

func unexportPin(baseDir string, pin int) error {
	return write(fmt.Sprintf("%s/unexport", baseDir),
		fmt.Sprintf("%d", pin))
}

func setDirection(baseDir string, pin int, direction int) error {
	pinDirection := map[int]string{1: "out", 0: "in"} // should/could be a const?

	return write(fmt.Sprintf("%s/gpio%d/direction", baseDir, pin),
		pinDirection[direction])
}

func writeValue(baseDir string, pin int, value int) error {
	return write(fmt.Sprintf("%s/gpio%d/value", baseDir, pin),
		fmt.Sprintf("%d", value))
}

func (g *RpiGPIO) Close() error {
	return unexportPin(g.baseDir, g.pin)
}

func (g *RpiGPIO) MakeOutput() error {
	if !g.isExported {
		err := exportPin(g.baseDir, g.pin)
		if err != nil {
			return attachErrorCause(fmt.Sprintf("Failed to export pin %d", g.pin), err)
		}
		g.isExported = true
	}

	if !g.isOutput {
		err := setDirection(g.baseDir, g.pin, 1)
		if err != nil {
			return attachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
		}

		g.isOutput = true
	}

	return nil
}

func (g *RpiGPIO) MakeInput() error {
	if !g.isExported {
		err := exportPin(g.baseDir, g.pin)
		if err != nil {
			return attachErrorCause(fmt.Sprintf("Failed to export pin %d", g.pin), err)
		}
		g.isExported = true
	}

	if g.isOutput {
		err := setDirection(g.baseDir, g.pin, 0)
		if err != nil {
			return attachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
		}
		g.isOutput = false
	}

	return nil
}

func (g *RpiGPIO) WriteValue(val int) error {
	if !g.isOutput {
		return &RpiGPIOError{msg: fmt.Sprintf("Pin %d is not an output pin", g.pin)}
	}

	err := writeValue(g.baseDir, g.pin, val)
	if err != nil {
		return attachErrorCause(fmt.Sprintf("Failed to set output to %d on pin %d", val, g.pin), err)
	}

	return nil
}

func (g *RpiGPIO) ReadValue() (int, error) {
	// return the value
	// /sys/class/gpio/gpio%d/value

	// skipping for a hot min... gonna see output work first
	return 0, nil
}

func NewRpiGPIO(pin int) (*RpiGPIO, error) {
	r := &RpiGPIO{pin: pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeOutput()
}

// for tests I can change the baseDir to some mock shit
// and read/cleanup files from there
// then in real life the GPIO iface would be mocked for tests
