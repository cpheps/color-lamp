//Package configloader responsible for loading configuration file
package configloader

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

const (
	configFileLocation = "./config/config.toml"
)

//LifeLineConfig represents "lifeline" part of config
type LifeLineConfig struct {
	ClusterName string `toml:"clusterName"`
	HostName    string `toml:"hostname"`
	Port        string `toml:"port"`
}

//LampConfig represents "lamp" part of config
type LampConfig struct {
	Green uint8 `toml:"green"`
	Red   uint8 `toml:"red"`
	Blue  uint8 `toml:"blue"`
}

//Config represents configuration file
type Config struct {
	LifeLineConfig LifeLineConfig `toml:"lifeline"`
	LampConfig     LampConfig     `toml:"lamp"`
}

//LoadConfig loads configuration file
func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		log.Printf("Error reading config file: %s", err.Error())
		return nil, err
	}

	config := &Config{}
	err = toml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error parsing config: %s", err.Error())
		return nil, err
	}

	return config, nil
}

//SaveConfig saves configuration file
func SaveConfig(config *Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		log.Printf("Error marshaling config: %s", err.Error())
		return err
	}

	err = ioutil.WriteFile(configFileLocation, data, os.ModePerm)
	if err != nil {
		log.Printf("Error saving config file: %s", err.Error())
		return err
	}

	return nil
}
