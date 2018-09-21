package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/cpheps/color-lamp/buttoncontroller"
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
	closeChan := make(chan bool)
	defer func() {
		closeChan <- true
		close(closeChan)
	}()

	eventChan := setupButton(closeChan)

	//Load in config
	config, err := configloader.LoadConfig()
	checkErr(err)

	//Init Lamp
	newLamp, err := setupLamp(&config.LampConfig)
	checkErr(err)
	defer newLamp.TearDown()

	//Init Client
	client := setupClient(&config.LifeLineConfig)

	//Core loop
	serverTicker := time.NewTicker(15 * time.Second).C

	//Signals if button press should override server read for next cycle.
	//Stops the case of pushing the button on this lamp and reading from the server
	//before the server is updated.
	override := false

	for {
		select {
		case <-serverTicker:
			//If we overriding skip checking the server
			if override {
				override = false
				continue
			}
			queryAndUpdate(client, newLamp, config.LifeLineConfig.ClusterID)
		case event := <-eventChan:
			if event == buttoncontroller.PressedEvent {
				client.SetClusterColor(config.LifeLineConfig.ClusterID, newLamp.GetLampColor())
				err := newLamp.SetCurrentColor(newLamp.GetLampColor())
				if err != nil {
					fmt.Println("Error setting Lamp color:", err.Error())
				} else {
					override = true
				}
				continue
			}

			//If not press then hold
			// run shutdown command in defer so lamp clean up happens.
			cmd := "sudo shutdown -h now"
			defer exec.Command("/bin/sh", "-c", cmd).Run()
			return
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
	color, err := client.GetClusterColor(clusterID)
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

func setupButton(closeChan <-chan bool) <-chan buttoncontroller.ButtonEvent {
	//Init buttons
	toggleButton, err := buttoncontroller.CreateButton(uint8(21))
	checkErr(err)

	return buttoncontroller.HandleButton(toggleButton, closeChan)
}

// TODO change this logic so it doesn't exit on a failure. Rather log
// all failures and continue on if possible
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
