package main

import (
	"fmt"
	"os"

	"github.com/cpheps/color-lamp/lamp"
	"github.com/cpheps/color-lamp/ledcontrol"
)

//Version current version of the Lamp Life Line server
var Version string

//BuildTime build time of binary
var BuildTime string

func main() {
	fmt.Printf("Running Color Lamp version %s build on %s\n", Version, BuildTime)

	lamp, err := setupLamp(uint32(0x200000))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer lamp.TearDown()
}

func setupLEDControl() (*ledcontrol.LEDControl, error) {
	ledControl, err := ledcontrol.CreateLEDControl(ledcontrol.DefaultPin, 16, ledcontrol.FullBrightness)
	if err != nil {
		return nil, fmt.Errorf("Error creating LEDControl: %s", err.Error())
	}

	err = ledControl.Init()
	if err != nil {
		return nil, fmt.Errorf("Error initializing LED Control: %s", err.Error())
	}

	return ledControl, nil
}

func setupLamp(lampColor uint32) (*lamp.Lamp, error) {
	ledControl, err := setupLEDControl()
	if err != nil {
		return nil, err
	}

	return lamp.CreateLamp(lampColor, lampColor, ledControl)
}
