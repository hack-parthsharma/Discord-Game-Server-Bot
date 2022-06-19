package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const filename = "config.yaml"

// Config is the bot's config settings
type Config struct {
	DiscordCfg Discord `yaml:"discord"`
}

// Discord is the Discord specific configuration
type Discord struct {
	Token    string   `yaml:"token"`
	Channels []string `yaml:"channels"`
}

// NewConfig returns a new decoded Config struct
func NewConfig() (*Config, error) {
	config := &Config{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading config file named: %v", filename)
	}

	if err = yaml.Unmarshal(file, config); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal yaml file")
	}

	return config, err
}
