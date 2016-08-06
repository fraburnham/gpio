package sysfs

import (
	"gpio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testBaseDir      = "."
	exportPathFmt    = "%s/export"
	unexportPathFmt  = "%s/unexport"
	directionPathFmt = "%s/gpio%d/direction"
	valuePathFmt     = "%s/gpio%d/value"
	passCheck        = "\u2713"
	failX            = "\u2717"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func newTestSysfsGPIO(pin int) *SysfsGPIO {
	return &SysfsGPIO{pin: pin,
		baseDir: testBaseDir}
}

func setUp(pin int) {
	files := []string{fmt.Sprintf(exportPathFmt, testBaseDir),
		fmt.Sprintf(unexportPathFmt, testBaseDir),
		fmt.Sprintf(directionPathFmt, testBaseDir, pin),
		fmt.Sprintf(valuePathFmt, testBaseDir, pin)}

	// why does this need to be 0777? Why isn't the directory actually 0777?
	// (its 0755)
	panicErr(os.Mkdir(fmt.Sprintf("%s/gpio%d/", testBaseDir, pin), 0777))

	for i := range files {
		_, err := os.Create(files[i])
		panicErr(err)
	}
}

func tearDown(pin int) {
	files := []string{fmt.Sprintf(exportPathFmt, testBaseDir),
		fmt.Sprintf(unexportPathFmt, testBaseDir),
		fmt.Sprintf(directionPathFmt, testBaseDir, pin),
		fmt.Sprintf(valuePathFmt, testBaseDir, pin)}

	for i := range files {
		panicErr(os.Remove(files[i]))
	}

	panicErr(os.Remove(fmt.Sprintf("%s/gpio%d/", testBaseDir, pin)))
}

func checkValue(filePath string, expectedValue string) bool {
	data, err := ioutil.ReadFile(filePath)
	panicErr(err)
	return string(data) == expectedValue
}

func TestMakeOutput(t *testing.T) {
	pinNum := 1
	pinStr := fmt.Sprintf("%d", pinNum)
	s := newTestSysfsGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)

	panicErr(s.MakeOutput())

	if !s.isOutput || !s.isExported {
		t.Errorf("%s MakeOutput failed to update SysfsGPIO", failX)
	}

	if !checkValue(fmt.Sprintf(exportPathFmt, testBaseDir), pinStr) {
		t.Errorf("%s MakeOutput failed to export pin", failX)
	}

	if !checkValue(fmt.Sprintf(directionPathFmt, testBaseDir, pinNum), "out") {
		t.Errorf("%s MakeOutput failed to set pin direction", failX)
	}

	tearDown(pinNum)
	t.Logf("%s MakeOutput", passCheck)
}

func TestMakeInput(t *testing.T) {
	pinNum := 1
	pinStr := fmt.Sprintf("%d", pinNum)
	s := newTestSysfsGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)

	panicErr(s.MakeInput())

	if s.isOutput || !s.isExported {
		t.Errorf("%s MakeOutput failed to update SysfsGPIO", failX)
	}

	if !checkValue(fmt.Sprintf(exportPathFmt, testBaseDir), pinStr) {
		t.Errorf("%s MakeInput failed to export pin", failX)
	}

	if !checkValue(fmt.Sprintf(directionPathFmt, testBaseDir, pinNum), "in") {
		t.Errorf("%s MakeInput failed to set pin direction", failX)
	}

	tearDown(pinNum)
	t.Logf("%s MakeInput", passCheck)
}

func TestWriteValue(t *testing.T) {
	pinNum := 1
	s := newTestSysfsGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)
	panicErr(s.MakeOutput())

	panicErr(s.WriteValue(1))

	if !checkValue(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "1") {
		t.Errorf("%s WriteValue failed to set pin high", failX)
	}

	panicErr(s.WriteValue(0))

	if !checkValue(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "0") {
		t.Errorf("%s WriteValue failed to set pin low", failX)
	}

	tearDown(pinNum)
	t.Logf("%s WriteValue", passCheck)
}

func TestReadValue(t *testing.T) {
	pinNum := 1
	s := newTestSysfsGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)
	s.MakeInput()

	panicErr(write(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "1"))

	val, err := s.ReadValue()
	panicErr(err)
	if val != 1 {
		t.Errorf("%s ReadValue failed to read pin set high", failX)
	}

	panicErr(write(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "0"))

	val, err = s.ReadValue()
	panicErr(err)
	if val != 0 {
		t.Errorf("%s ReadValue failed to read pin set low", failX)
	}

	tearDown(pinNum)
	t.Logf("%s ReadValue", passCheck)
}

func TestInterfaceImplementation(t *testing.T) {
	var x gpio.GPIO = newTestSysfsGPIO(1)
	defer x.Close()
	if x == nil {
		t.Errorf("SysfsGPIO does not implement GPIO!")
	}
}

func ExampleNewSysfsOutput() {
	// create a new output pin
	sysfsOutput, err := NewSysfsOutput(1)
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
	sysfsInput, err := NewSysfsInput(1)
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
