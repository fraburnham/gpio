package sysfs

import (
	"gpio"
	"fmt"
	"strconv"
)

type SysfsGPIO struct {
	pin        int
	isOutput   bool
	isExported bool
	baseDir    string
}

func (g *SysfsGPIO) Close() error {
	return unexportPin(g.baseDir, g.pin)
}

func (g *SysfsGPIO) MakeOutput() error {
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

	g.isOutput = true

	return nil
}

func (g *SysfsGPIO) MakeInput() error {
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
	g.isOutput = false

	return nil
}

func (g *SysfsGPIO) WriteValue(val int) error {
	if !g.isOutput {
		return gpio.NewGPIOError(fmt.Sprintf("Pin %d is not an output pin", g.pin))
	}

	err := writeValue(g.baseDir, g.pin, val)
	if err != nil {
		return gpio.AttachErrorCause(fmt.Sprintf("Failed to set output to %d on pin %d", val, g.pin), err)
	}

	return nil
}

func (g *SysfsGPIO) ReadValue() (int, error) {
	if g.isOutput {
		return 0, gpio.NewGPIOError(fmt.Sprintf("Pin %d is not an input pin", g.pin))
	}

	data, err := readValue(g.baseDir, g.pin)

	if err != nil {
		return 0, gpio.AttachErrorCause(fmt.Sprintf("Failed to read value from pin %d", g.pin), err)
	}

	return strconv.Atoi(data)
}

func NewSysfsOutput(pin int) (*SysfsGPIO, error) {
	r := &SysfsGPIO{pin: pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeOutput()
}

func NewSysfsInput(pin int) (*SysfsGPIO, error) {
	r := &SysfsGPIO{pin: pin,
		baseDir: "/sys/class/gpio"}

	return r, r.MakeInput()
}
