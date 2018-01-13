package main

import (
	"fmt"
	"os"

	"github.com/cpheps/color-lamp/lampclient"

	"github.com/cpheps/color-lamp/configloader"
	"github.com/cpheps/color-lamp/lamp"
	"github.com/cpheps/color-lamp/ledcontrol"
)

//Version current version of the Lamp Life Line server
var Version string

//BuildTime build time of binary
var BuildTime string

func main() {
	fmt.Printf("Running Color Lamp version %s build on %s\n", Version, BuildTime)
	config, err := configloader.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	lamp, err = setupLamp(config.LampConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	client := setupClient(config.LifeLineConfig)

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

func setupLamp(config *configloader.LampConfig) (*lamp.Lamp, error) {
	ledControl, err := setupLEDControl()
	if err != nil {
		return nil, err
	}

	lampColor := lamp.ConvertRGB(config.Red, config.Green, config.Blue)
	return lamp.CreateLamp(lampColor, lampColor, ledControl)
}

func setupClient(config *configloader.LifeLineConfig) *lampclient.LampClient {
	return lampclient.CreateLampClient(config.HostName, config.Port)
}
