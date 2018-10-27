package buttoncontroller

import (
	"sync"
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
func HandleButton(b *Button, closeChan <-chan bool, wg *sync.WaitGroup) <-chan ButtonEvent {
	eventChan := make(chan ButtonEvent, 1)

	wg.Add(1)
	go buttonLoop(b, eventChan, closeChan, wg)

	return eventChan
}

func buttonLoop(b *Button, eventChan chan ButtonEvent, closeChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	buttonTicker := time.NewTicker(100 * time.Millisecond)
	var pressStart time.Time

	for {
		select {
		case closing := <-closeChan:
			if closing {
				buttonTicker.Stop()
				TearDown()
				return
			}
		case <-buttonTicker.C:
			switch b.buttonState {
			case idle:
				if b.IsPressed() {
					b.buttonState = pressed
					pressStart = time.Now()
				}
			case pressed:
				if b.IsPressed() {
					if secondsHeld := time.Since(pressStart); secondsHeld.Seconds() >= 5 {
						b.buttonState = coolDown
						eventChan <- HoldEvent
					}
				} else {
					b.buttonState = coolDown
					eventChan <- PressedEvent
				}
			case coolDown:
				<-time.NewTimer(5 * time.Second).C
				b.buttonState = idle
			}
		}
	}
}
