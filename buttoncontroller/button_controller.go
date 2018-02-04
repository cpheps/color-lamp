package buttoncontroller

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

var isInitialized = false

//Button represents a single button on a pin
type Button struct {
	pin rpio.Pin
}

//CreateButton Creates a new putton on the given pin
func CreateButton(pinNumber uint8) (*Button, error) {
	if !isInitialized {
		if err := rpio.Open(); err != nil {
			return nil, fmt.Errorf("Unable to initialize GPIO: %s", err.Error())
		}

		isInitialized = true
	}

	pin := rpio.Pin(pinNumber)
	pin.Input()
	pin.PullUp()

	return &Button{
		pin: rpio.Pin(pinNumber),
	}, nil
}

//IsPressed gets if the button is currently pressed
func (b Button) IsPressed() bool {
	//Since pull state is up a button press will give a state of low
	state := b.pin.Read()
	fmt.Println("Button:", state)
	return state == rpio.Low
}

//TearDown will deinit the button controller
func TearDown() error {
	return rpio.Close()
}