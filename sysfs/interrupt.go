package sysfs

/*
#include <fcntl.h>
#include <poll.h>
#include <stdlib.h>
#include <stdio.h>

// stolen from https://developer.ridgerun.com/wiki/index.php?title=Gpio-int-test.c
#define MAX_BUF 64

struct interrupt {
  int result;
  int value;
};

// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt

// when timeout_ms is -1 wait forever
// returns 1 on an event 0 on timeout -1 on error (like poll)
struct interrupt interrupt_poll(int fd, int timeout_ms) {
  struct pollfd fds[1];
  struct interrupt i[1];

  fds[0].fd = fd;
  fds[0].events = POLLPRI;

  i[0].result = poll(fds, 1, timeout_ms);
  if (fds[0].revents & POLLPRI) {
    char buf[MAX_BUF];
    lseek(fd, 0, SEEK_SET);
    read(fd, &buf, MAX_BUF);
    i[0].value = atoi(buf);
    printf("\nPoll result: %d, Read val: %d\n", i[0].result, i[0].value); // leaving this in until rangefinder is fixe
    // also this will need to make it to github
  }

  return i[0];
}

int c_read(int fd) {
  char *buf[MAX_BUF];
  read(fd, buf, MAX_BUF);
}
*/
import "C"

import (
	"fmt"
	"github.com/fraburnham/gpio"
	"os"
	"time"
)

func validEdge(edge string) bool {
	validEdges := map[string]bool{
		"rising":  true,
		"falling": true,
		"both":    true}
	return validEdges[edge]
}

func interruptListner(ctrlCh chan bool, g *GPIO, ch chan gpio.InterruptEvent, pollTimeoutMs int) {
	f, err := os.OpenFile(valueFmt(g.baseDir, g.pin), os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		ch <- gpio.InterruptEvent{Err: err}
		g.wg.Done()
		return
	}

	fd := C.int(f.Fd())
	C.c_read(fd) // not sure if this needs to be read from c, should experiment

	for { // it would be nice to refactor some of this, the conditionals are getting hard to read
		select {
		case <-ctrlCh:
			g.wg.Done()
			return
		default:
			poll := C.interrupt_poll(fd, C.int(pollTimeoutMs))
			switch poll.result {
			case 1:
				eventTime := time.Now()
				ch <- gpio.InterruptEvent{
					Value:     int(poll.value),
					Timestamp: eventTime}
			case -1:
				ch <- gpio.InterruptEvent{
					Err: gpio.NewError(fmt.Sprintf("poll failed on interrupt for pin %d", g.pin))}
			case 0:
				break
			}
		}
	}
}

// SetInterrupt sets a rising, falling or both type interrupt on a gpio. InterruptEvents are
// placed onto a chan for consumption. The pollTimeoutMs determines how long the interrupt
// routine should block before checking if the interrupt has been cleared. If you need to
// quickly clear and interrupt and use the gpio for something else this value should be fairly
// low. However a lower value will result in more CPU usage.
func (g *GPIO) SetInterrupt(edge string, ch chan gpio.InterruptEvent, pollTimeoutMs int) error {
	if g.isInterrupt {
		return gpio.NewError(fmt.Sprintf("Pin %d already has an interrupt set", g.pin))
	}

	if !g.isExported {
		return gpio.NewError(fmt.Sprintf("Pin %d is not exported", g.pin))
	}

	if g.direction == "output" {
		return gpio.NewError(fmt.Sprintf("Pin %d is not an input pin", g.pin))
	}

	if !validEdge(edge) {
		return gpio.NewError(fmt.Sprintf("Unable to set edge. Got %s expected one of rising, falling, both.", edge))
	}

	err := setInterrupt(g.baseDir, g.pin, edge)
	if err != nil {
		return gpio.AttachErrorCause(fmt.Sprintf("Failed to set interrupt on pin %d with edge %s", g.pin, edge), err)
	}

	g.interruptCtrlCh = make(chan bool)
	g.wg.Add(1)
	g.isInterrupt = true
	go interruptListner(g.interruptCtrlCh, g, ch, pollTimeoutMs)

	return nil
}

// ClearInterrupt clears an interrupt on a gpio.
func (g *GPIO) ClearInterrupt() error {
	if !g.isInterrupt {
		return gpio.NewError(fmt.Sprintf("Pin %d does not have an interrupt set", g.pin))
	}

	g.interruptCtrlCh <- true
	return nil
}
