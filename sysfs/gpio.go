package sysfs

import (
	"fmt"
	"gpio"
	"strconv"
	"sync"
)

type GPIO struct {
	pin             int
	direction       string
	isExported      bool
	isInterrupt     bool
	interruptCtrlCh chan bool
	baseDir         string
	wg              sync.WaitGroup
}

func (g *GPIO) Close() error {
	g.wg.Wait()
	return unexportPin(g.baseDir, g.pin)
}

func (g *GPIO) MakeOutput() error {
	if !g.isExported {
		err := exportPin(g.baseDir, g.pin)
		if err != nil {
			return gpio.AttachErrorCause(fmt.Sprintf("Failed to export pin %d", g.pin), err)
		}
		g.isExported = true
	}

	err := setDirection(g.baseDir, g.pin, 1)
	if err != nil {
		return gpio.AttachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
	}

	g.direction = "output"

	return nil
}

func (g *GPIO) MakeInput() error {
	if !g.isExported {
		err := exportPin(g.baseDir, g.pin)
		if err != nil {
			return gpio.AttachErrorCause(fmt.Sprintf("Failed to export pin %d", g.pin), err)
		}
		g.isExported = true
	}

	err := setDirection(g.baseDir, g.pin, 0)
	if err != nil {
		return gpio.AttachErrorCause(fmt.Sprintf("Failed to set pin %d direction", g.pin), err)
	}
	g.direction = "input"

	return nil
}

func (g *GPIO) WriteValue(val int) error {
	if !g.isExported {
		return gpio.NewError(fmt.Sprintf("Pin %d is not exported", g.pin))
	}

	if g.direction == "input" {
		return gpio.NewError(fmt.Sprintf("Pin %d is not an output pin", g.pin))
	}

	err := writeValue(g.baseDir, g.pin, val)
	if err != nil {
		return gpio.AttachErrorCause(fmt.Sprintf("Failed to set output to %d on pin %d", val, g.pin), err)
	}

	return nil
}

func (g *GPIO) ReadValue() (int, error) {
	if !g.isExported {
		return 0, gpio.NewError(fmt.Sprintf("Pin %d is not exported", g.pin))
	}

	if g.direction == "output" {
		return 0, gpio.NewError(fmt.Sprintf("Pin %d is not an input pin", g.pin))
	}

	data, err := readValue(g.baseDir, g.pin)

	if err != nil {
		return 0, gpio.AttachErrorCause(fmt.Sprintf("Failed to read value from pin %d", g.pin), err)
	}

	return strconv.Atoi(data)
}

func NewOutput(pin int) (*GPIO, error) {
	r := &GPIO{
		pin:     pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeOutput()
}

func NewInput(pin int) (*GPIO, error) {
	r := &GPIO{
		pin:     pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeInput()
}
