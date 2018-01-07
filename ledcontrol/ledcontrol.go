package ledcontrol

import (
	"errors"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

const (
	//DefaultPin is the default GPIO pin on a raspberry pi
	DefaultPin = 18

	//FullBrightness is the maximum brightness for the LEDs
	FullBrightness = 255
)

//LEDControl represents the control module for the LED lights
type LEDControl struct {
	pin        int
	ledCount   int
	brightness int

	init bool
}

//CreateLEDControl creates a new instance of LEDControl
func CreateLEDControl(pin, ledCount, brightness int) (*LEDControl, error) {
	if brightness < 0 || brightness > 255 {
		return nil, errors.New("Brightness must be between 0 and 255")
	}

	return &LEDControl{
		pin:        pin,
		ledCount:   ledCount,
		brightness: brightness,
	}, nil
}

//Init initializes the LED Strip
func (lc *LEDControl) Init() error {
	if lc.init {
		return errors.New("LED Control already initialized")
	}

	err := ws2811.Init(lc.pin, lc.ledCount, lc.brightness)
	if err != nil {
		return err
	}

	lc.init = true
	return nil
}

//DeInit de-initializes the LED Strip after use
func (lc *LEDControl) DeInit() error {
	if !lc.init {
		return errors.New("LED Control has not been initialized yet")
	}

	ws2811.Clear()
	ws2811.Fini()
	lc.init = false

	return nil
}

//ChangeStripColor changes the color of the LED Strip
func (lc LEDControl) ChangeStripColor(color uint32) error {
	if !lc.init {
		return errors.New("LED Control has not been initialized yet")
	}

	for i := 0; i < lc.ledCount; i++ {
		ws2811.SetLed(i, color)
	}

	return ws2811.Render()
}
