package events

import (
	"gpio"
	"time"
)

type EdgeEvent struct {
	BeforeEvent int
	AfterEvent int
	Timestamp time.Time
}

func edgeTrigger(pin gpio.GPIO, eventCh chan EdgeEvent, ctrlCh chan bool) (error) {
	lastState, err := pin.ReadValue()
	if err != nil {
		panic(err) // improve
	}

	for true {
		select {
		case <-ctrlCh:
			return nil
		default:
			newState, err := pin.ReadValue()
			if err != nil {
				panic(err)  // improve
			}

			if newState != lastState {
				eventCh <- EdgeEvent{BeforeEvent: lastState,
					AfterEvent: newState,
					Timestamp: time.Now()}
				lastState = newState
			}
		}
	}
	return nil
}

func StartEdgeTrigger(pin gpio.GPIO, holdEvents int) (chan EdgeEvent, chan bool) {
	eventCh := make(chan EdgeEvent, holdEvents) // this should have a buffer
	ctrlCh := make(chan bool)

	go edgeTrigger(pin, eventCh, ctrlCh)

	return eventCh, ctrlCh
}

func StopEdgeTrigger(ctrlCh chan bool) {
	ctrlCh <- true
}
