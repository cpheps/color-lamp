package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cpheps/color-lamp/configloader"
	"github.com/cpheps/color-lamp/lamp"
	"github.com/cpheps/color-lamp/lampclient"
	"github.com/cpheps/color-lamp/ledcontrol"
)

//Version current version of the Lamp Life Line server
var Version string

//BuildTime build time of binary
var BuildTime string

func main() {
	fmt.Printf("Running Color Lamp version %s build on %s\n", Version, BuildTime)

	//Load in config
	config, err := configloader.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//Init Lamp
	lamp, err = setupLamp(config.LampConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer lamp.TearDown()

	//Init Client
	client := setupClient(config.LifeLineConfig)

	//Core loop
	serverTicker := time.NewTicker(15 * time.Second).C

	for {
		select {
		case <-serverTicker:
			queryAndUpdate(client, lamp, config.LifeLineConfig.ClusterID)
		}
	}
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

func queryAndUpdate(client *lampclient.LampClient, lamp *lamp.Lamp, clusterID string) {
	color, err := client.GetClusterColor()
	if err != nil {
		fmt.Println("Error getting cluster color:", err.Error())
		return
	}

	if *color != lamp.GetCurrentColor() {
		err = lamp.SetCurrentColor(*color)
		if err != nil {
			fmt.Println("Error setting Lamp color:", err.Error())
			return
		}
	}
	fmt.Println("Set Lamp color to:", *color)
}
