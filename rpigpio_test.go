package gpio

import (
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

func newTestRpiGPIO(pin int) *RpiGPIO {
	return &RpiGPIO{pin: pin,
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
	r := newTestRpiGPIO(pinNum)
	defer r.Close()
	setUp(pinNum)

	panicErr(r.MakeOutput())

	if !r.isOutput || !r.isExported {
		t.Errorf("%s MakeOutput failed to update RpiGPIO", failX)
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
	r := newTestRpiGPIO(pinNum)
	defer r.Close()
	setUp(pinNum)

	panicErr(r.MakeInput())

	if r.isOutput || !r.isExported {
		t.Errorf("%s MakeOutput failed to update RpiGPIO", failX)
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
	r := newTestRpiGPIO(pinNum)
	defer r.Close()
	setUp(pinNum)
	panicErr(r.MakeOutput())

	panicErr(r.WriteValue(1))

	if !checkValue(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "1") {
		t.Errorf("%s WriteValue failed to set pin high", failX)
	}

	panicErr(r.WriteValue(0))

	if !checkValue(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "0") {
		t.Errorf("%s WriteValue failed to set pin low", failX)
	}

	tearDown(pinNum)
	t.Logf("%s WriteValue", passCheck)
}

func TestReadValue(t *testing.T) {
	pinNum := 1
	r := newTestRpiGPIO(pinNum)
	defer r.Close()
	setUp(pinNum)
	r.MakeInput()

	panicErr(write(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "1"))

	val, err := r.ReadValue()
	panicErr(err)
	if val != 1 {
		t.Errorf("%s ReadValue failed to read pin set high", failX)
	}

	panicErr(write(fmt.Sprintf(valuePathFmt, testBaseDir, pinNum), "0"))

	val, err = r.ReadValue()
	panicErr(err)
	if val != 0 {
		t.Errorf("%s ReadValue failed to read pin set low", failX)
	}

	tearDown(pinNum)
	t.Logf("%s ReadValue", passCheck)
}

func ExampleNewRpiOutput() {
	// create a new output pin
	rpiOutput, err := NewRpiOutput(1)
	if err != nil {
		panic(err)
	}

	// set pin high
	err = rpiOutput.WriteValue(1)
	if err != nil {
		panic(err)
	}

	// set pin low
	err = rpiOutput.WriteValue(0)
	if err != nil {
		panic(err)
	}
}

func ExampleNewRpiInput() {
	// create new input pin
	rpiInput, err := NewRpiInput(1)
	if err != nil {
		panic(err)
	}

	// read value
	val, err := rpiInput.ReadValue()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
