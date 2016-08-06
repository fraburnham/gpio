package sysfs

/*
#include <fcntl.h>
#include <poll.h>
#include <stdlib.h>
#include <stdio.h>

// stolen from https://developer.ridgerun.com/wiki/index.php?title=Gpio-int-test.c
#define MAX_BUF 64

// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt

// when timeout_ms is -1 wait forever
// returns 1 on an event 0 on timeout -1 on error (like poll)
int interrupt_poll(int fd, int timeout_ms) {
  struct pollfd fds[1];

  fds[0].fd = fd;
  fds[0].events = POLLPRI;

  int result = poll(fds, 1, timeout_ms);
  if (fds[0].revents & POLLPRI) {
    lseek(fd, 0, SEEK_SET);
    char *buf[MAX_BUF];
    read(fd, buf, MAX_BUF);
  }

  return result;
}

int c_read(int fd) {
  char *buf[MAX_BUF];
  read(fd, buf, MAX_BUF);
}
*/
import "C"

import (
	"fmt"
	"gpio"
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

func interruptListner(ctrlCh chan bool, g *GPIO, ch chan gpio.InterruptEvent) {
	f, err := os.OpenFile(valueFmt(g.baseDir, g.pin), os.O_RDONLY, 0644)
	fd := C.int(f.Fd())
	C.c_read(fd)
	if err != nil {
		panic(err) // shitty!
	}

	defer f.Close()

	for {
		select {
		case <-ctrlCh:
			g.wg.Done()
			return
		default:
			poll := C.interrupt_poll(fd, 10) // this polling frequency should be configurable!
			switch poll {
			case 1:
				eventTime := time.Now()
				val, _ := g.ReadValue()          // needs error handling
				ch <- gpio.InterruptEvent{
					Value: val,
					Timestamp: eventTime}
			// case -1: some kind of error handling
			case 0:
				break
			}
		}
	}
}

func (g *GPIO) SetInterrupt(edge string, ch chan gpio.InterruptEvent) (chan bool, error) {
	if g.isInterrupt {
		return nil, gpio.NewGPIOError(fmt.Sprintf("Pin %d already has an interrupt set", g.pin))
	}

	if !g.isExported {
		return nil, gpio.NewGPIOError(fmt.Sprintf("Pin %d is not exported", g.pin))
	}

	if g.isOutput { // this should be updated to check direction (less error prone)
		return nil, gpio.NewGPIOError(fmt.Sprintf("Pin %d is not an input pin", g.pin))
	}

	if !validEdge(edge) {
		return nil, gpio.NewGPIOError(fmt.Sprintf("Unable to set edge. Got %s expected one of rising, falling, both.", edge))
	}

	err := setInterrupt(g.baseDir, g.pin, edge)
	if err != nil {
		return nil, gpio.AttachErrorCause(fmt.Sprintf("Failed to set interrupt on pin %d with edge %s", g.pin, edge), err)
	}

	ctrlCh := make(chan bool)
	g.wg.Add(1)
	g.isInterrupt = true
	go interruptListner(ctrlCh, g, ch)

	return ctrlCh, nil
}

func (g *GPIO) ClearInterrupt(ctrlCh chan bool) error {
	if !g.isInterrupt {
		// rename to gpio.NewError
		return gpio.NewGPIOError(fmt.Sprintf("Pin %d does not have an interrupt set", g.pin))
	}

	ctrlCh <- true
	return nil
}
