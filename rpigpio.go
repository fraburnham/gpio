package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type RpiGPIO struct {
	pin        int
	isOutput   bool
	isExported bool
	baseDir    string
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

	err := setDirection(g.baseDir, g.pin, 1)
	if err != nil {
		return attachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
	}

	g.isOutput = true

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

	err := setDirection(g.baseDir, g.pin, 0)
	if err != nil {
		return attachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
	}
	g.isOutput = false

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
	if g.isOutput {
		return 0, &RpiGPIOError{msg: fmt.Sprintf("Pin %d is not an input pin", g.pin)}
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/gpio%d/value", g.baseDir, g.pin))

	if err != nil {
		return 0, attachErrorCause(fmt.Sprintf("Failed to read value from pin %d", g.pin), err)
	}

	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func NewRpiOutput(pin int) (*RpiGPIO, error) {
	r := &RpiGPIO{pin: pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeOutput()
}

func NewRpiInput(pin int) (*RpiGPIO, error) {
	r := &RpiGPIO{pin: pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeInput()
}
