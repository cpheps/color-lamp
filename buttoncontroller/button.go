package buttoncontroller

import (
	"log"

	"github.com/stianeikeland/go-rpio"
)

type state int

const (
	idle state = iota
	pressed
	coolDown
)

var isInitialized = false

//Button represents a single button on a pin
type Button struct {
	pin         rpio.Pin
	buttonState state
}

//CreateButton Creates a new putton on the given pin
func CreateButton(pinNumber uint8) (*Button, error) {
	if !isInitialized {
		if err := rpio.Open(); err != nil {
			log.Printf("unable to initialize GPIO: %s", err.Error())
			return nil, err
		}

		isInitialized = true
	}

	pin := rpio.Pin(pinNumber)
	pin.Input()
	pin.PullUp()

	return &Button{
		pin:         rpio.Pin(pinNumber),
		buttonState: idle,
	}, nil
}

//IsPressed gets if the button is currently pressed
func (b Button) IsPressed() bool {
	//Since pull state is up a button press will give a state of low
	return b.pin.Read() == rpio.Low
}

//TearDown will deinit the button controller
func TearDown() error {
	return rpio.Close()
}
