package lamp

import (
	"errors"
	"fmt"

	"github.com/cpheps/color-lamp/ledcontrol"
)

//Lamp represents a single lamp object
type Lamp struct {
	lampColor    uint32
	currentColor uint32

	ledControl *ledcontrol.LEDControl
}

//CreateLamp creates a new lamp object
//Assumes LEDControl Init is run before passing in
func CreateLamp(lampColor, currentColor uint32, ledControl *ledcontrol.LEDControl) (*Lamp, error) {
	if ledControl == nil {
		return nil, errors.New("Error creating Lamp, ledControl must not be nil")
	}

	lamp := &Lamp{
		lampColor:  lampColor,
		ledControl: ledControl,
	}

	err := lamp.SetCurrentColor(currentColor)
	if err != nil {
		return nil, fmt.Errorf("Error setting Lamp, color on startup: %s", err.Error())
	}

	return lamp, nil
}

//GetCurrentColor returns the current color
func (l Lamp) GetCurrentColor() uint32 {
	return l.currentColor
}

//SetLampColor sets the color the belongs to this lamp
func (l *Lamp) SetLampColor(red uint8, green uint8, blue uint8) {
	l.lampColor = ConvertRGB(red, green, blue)
}

//SetCurrentColor sets the current color of the lamp
func (l *Lamp) SetCurrentColor(color uint32) error {
	err := l.ledControl.ChangeStripColor(color)
	if err == nil {
		l.currentColor = color
	}

	return err
}

//TearDown Deinits LEDControl
func (l Lamp) TearDown() {
	l.ledControl.Deinit()
}

//ConvertRGB Converts and RGB color to uint32
func ConvertRGB(red uint8, green uint8, blue uint8) uint32 {
	return uint32(green)<<16 | uint32(red)<<8 | uint32(blue)
}
