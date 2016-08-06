package sysfs

import (
	"fmt"
	"gpio"
	"testing"
)

func setValue(g *GPIO, value string) {
	panicErr(write(fmt.Sprintf(valuePathFmt, g.baseDir, g.pin), value))
}

func TestUnexportedPin(t *testing.T) {
	pinNum := 1
	s := newTestGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)

	_, err := s.SetInterrupt("rising", make(chan gpio.InterruptEvent))
	if err == nil {
		t.Errorf("%s interruptListenr didn't fail on an unexported pin", failX)
	}

	tearDown(pinNum)
	t.Logf("%s Interrupt unexported pin error", passCheck)
}

/*
The sysfs files aren't regular files. I'm not yet sure how to mock them
for testing. `poll`ing regular files has different results and doesn't work
like an interrupt

func TestInterruptListenerBothEdge(t *testing.T) {
	pinNum := 1
	s := newTestGPIO(pinNum)
	defer s.Close()
	setUp(pinNum)
	s.MakeInput()

	// some real tests:
	// set the listener (edge both) and write a 1 to the value file check that a message is published within 5ms
	// (again the poll timeout should be an option)
	// with the same listener write a 0 to the file and see that another message was published
	// close the listener
	// clean value file

	ch := make(chan gpio.InterruptEvent)
	//defer close(ch)
	ctrlCh, err := s.Interrupt("both", ch)
	panicErr(err) // this is a shitty way to do this failure, it works for now
	// need a fn to write the values
	setValue(s, "1")
	time.Sleep(time.Duration(10)*time.Millisecond)
	setValue(s, "0")
	time.Sleep(time.Duration(10)*time.Millisecond)
	setValue(s, "1")
	time.Sleep(time.Duration(10)*time.Millisecond)
	t.Logf("Back from sleep\n")
	t.Logf("%s\n", checkValue(fmt.Sprintf(valuePathFmt, s.baseDir, s.pin), "1"))
	select {
	case val := <- ch:
		t.Logf("%s\n", val)
		t.Logf("%s Interrupt 'both' edge sees rising edge", passCheck)
	default:
		t.Errorf("%s Interrupt failed to see rising change within 5ms", failX)
	}

	// set listener (edge rising) and write a 1 to the value file
	// see message
	// write 0, see no message

	// set listener (edge falling) and write 1
	// see no message
	// write 0
	// see message

	// make an output pin
	// see fail

	// test pass

	s.ClearInterrupt(ctrlCh)
	tearDown(pinNum)
	//t.Logf("%s Interrupt 'both' edge", passCheck)
}
*/
