package main

import (
	"fmt"
	"os"
	"time"

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

	defer ledControl.DeInit()

	_, err = lamp.CreateLamp(uint32(0x200000), uint32(0x200000), ledControl)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	time.Sleep(1 * time.Minute)
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
