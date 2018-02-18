package buttoncontroller

import (
	"fmt"
	"time"
)

type state int

const (
	idle state = iota
	pressed
	hold
)

//HandleButton controls button IO
//meant to be run as separate goroutine
func HandleButton(b *Button, closeChan <-chan bool) {
	buttonTicker := time.NewTicker(100 * time.Millisecond)
	var pressStart time.Time

	buttonState := idle

	for {
		select {
		case close := <-closeChan:
			if close {
				buttonTicker.Stop()
				return
			}
		case <-buttonTicker.C:
			switch buttonState {
			case idle:
				if b.IsPressed() {
					buttonState = pressed
					pressStart = time.Now()
				}
			case pressed:
				if b.IsPressed() {
					if secondsHeld := time.Now().Sub(pressStart) * time.Second; secondsHeld >= 3 {
						buttonState = hold
						fmt.Println("Button Hold")
					}
				} else {
					buttonState = idle
					fmt.Println("Button pressed")
				}
			case hold:
				if !b.IsPressed() {
					buttonState = idle
				}
			}
		}
	}
}
