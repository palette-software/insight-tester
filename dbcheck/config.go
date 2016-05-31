package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ConstParams struct {
	Port     int    `yaml:"Port"`
	Database string `yaml:"Database"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

type Databases struct {
	Hosts  []string    `yaml:"Hosts"`
	Params ConstParams `yaml:"Params"`
}

type Config struct {
	Databases           Databases `yaml:"Databases"`
	SplunkServerAddress string    `yaml:"SplunkServerAddress"`
	SplunkToken         string    `yaml:"SplunkToken"`
}

func parseConfig(fileName string) (*Config, error) {
	var v Config
	input, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer input.Close()
	b, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	if len(v.Databases.Hosts) < 1 ||
		v.Databases.Params.Port == 0 ||
		v.Databases.Params.User == "" ||
		v.Databases.Params.Database == "" ||
		v.Databases.Params.Password == "" {
		return nil, fmt.Errorf("Config file does not contain database information.")
	}
	return &v, nil
}
