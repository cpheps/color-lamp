package buttoncontroller

import (
	"time"
)

//ButtonEvent is a type of event when a button is pressed
type ButtonEvent int

const (
	//PressedEvent signals a button press
	PressedEvent ButtonEvent = iota

	//HoldEvent signals a hold of over 5 seconds
	HoldEvent
)

//HandleButton controls button IO
//
//Starts a separate goroutine to handle button IO.
//Communicates over channel
//On closing will handle button TearDown
func HandleButton(b *Button, closeChan <-chan bool) <-chan ButtonEvent {
	eventChan := make(chan ButtonEvent)

	go buttonLoop(b, eventChan, closeChan)

	return eventChan
}

func buttonLoop(b *Button, eventChan chan ButtonEvent, closeChan <-chan bool) {
	buttonTicker := time.NewTicker(100 * time.Millisecond)
	var pressStart time.Time

	buttonState := idle

	for {
		select {
		case closing := <-closeChan:
			if closing {
				buttonTicker.Stop()
				TearDown()
				close(eventChan)
				return
			}
		case <-buttonTicker.C:
			switch buttonState {
			case idle:
				if b.IsPressed() {
					b.buttonState = pressed
					pressStart = time.Now()
				}
			case pressed:
				if b.IsPressed() {
					if secondsHeld := time.Now().Sub(pressStart) * time.Second; secondsHeld >= (5 * time.Second) {
						b.buttonState = coolDown
						eventChan <- HoldEvent
					}
				} else {
					b.buttonState = coolDown
					eventChan <- PressedEvent
				}
			case coolDown:
				<-time.NewTimer(3 * time.Second).C
				b.buttonState = idle
			}
		}
	}
}
