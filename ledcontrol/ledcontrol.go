package ledcontrol

import (
	"image/color"

	"github.com/mcuadros/go-rpi-ws281x"
)

const (
	height = 1
	width = 4
)

type LEDControl struct {
	ledColor color.RGBA
	// canvase 
}

func CreateLEDControl() *LEDControl {


	return &LEDControl {
		// ledColor
	}
}