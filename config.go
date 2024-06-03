package main

import (
	"flag"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var input string

func init() {
	flag.StringVar(&input, "config", "/etc/convert/config.yaml", "/etc/convert/config.yaml")
}

type Config struct {
	Addr          string `yaml:"Addr"`
	OutputDir     string `yaml:"OutputDir"`
	GotenbergAddr string `yaml:"GotenbergAddr"`
}

func InitConfig() (*Config, error) {
	bytes, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, err
	}

	return config, nil
}
