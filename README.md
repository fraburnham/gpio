# gpio

The gpio package aims to provide a simple, reusable interface to enable the creation of software on boards like the raspberry pi and beaglebone.

Unless otherwise noted, the gpio source files are distributed
under the BSD-style license found in the LICENSE file.

--

## Quickstart

```go
func ExampleNewSysfsOutput() {
	// create a new output pin
	sysfsOutput, err := NewOutput(1)
	if err != nil {
		panic(err)
	}

	// set pin high
	err = sysfsOutput.WriteValue(1)
	if err != nil {
		panic(err)
	}

	// set pin low
	err = sysfsOutput.WriteValue(0)
	if err != nil {
		panic(err)
	}
}

func ExampleNewSysfsInput() {
	// create new input pin
	sysfsInput, err := NewInput(1)
	if err != nil {
		panic(err)
	}

	// read value
	val, err := sysfsInput.ReadValue()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
```

--

## API

See the [godocs](https://godoc.org/github.com/fraburnham/gpio) for details.
