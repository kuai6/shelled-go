package main

import (
	"fmt"

	"github.com/creasty/defaults"
)

type Config struct {
	Http struct {
		Host string `yaml:"host" default:"0.0.0.0"`
		Port int    `yaml:"port" default:"8080"`
	} `yaml:"http"`
	UDP struct {
		Host string `yaml:"host" default:"0.0.0.0"`
		Port int    `yaml:"port" default:"8830"`
	} `yaml:"udp"`
}

func LoadConfig(file string) (*Config, error) {

	config := &Config{}
	if err := defaults.Set(config); err != nil {
		return nil, fmt.Errorf("load config error: load defaults: %s", err)
	}

	return config, nil
}
