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

	ledControl, err := setupLEDControl()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	lamp, err := lamp.CreateLamp(uint32(0x020000), uint32(0x020000), ledControl)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setupLEDControl() (*ledcontrol.LEDControl, error) {
	ledControl, err := ledcontrol.CreateLEDControl()
	if err != nil {
		return nil, fmt.Errorf("Error creating LEDControl: %s", err.Error())
	}

	err = ledControl.Init()
	if err != nil {
		return nil, fmt.Errorf("Error initializing LED Control: %s", err.Error())
	}

	return ledControl, nil
}
