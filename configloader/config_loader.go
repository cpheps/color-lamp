//Package configloader responsible for loading configuration file
package configloader

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

const (
	configFileLocation = "./config/config.toml"
)

//LifeLineConfig represents "lifeline" part of config
type LifeLineConfig struct {
	HostName string `toml:"hostname"`
	Port     string `toml:"port"`
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
		return nil, fmt.Errorf("Error reading config file: %s", err.Error())
	}

	config := &Config{}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config: %s", err.Error())
	}

	return config, nil
}
