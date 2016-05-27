package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Database struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Database string `yaml:"Database"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

type Config struct {
	Database Database `yaml:"Database"`
}

func getResultDBConfig(fileName string) (*Database, error) {
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
	if v.Database.Host == "" ||
		v.Database.Port == 0 ||
		v.Database.User == "" ||
		v.Database.Database == "" ||
		v.Database.Password == "" {
		return nil, fmt.Errorf("Config file does not contain database information.")
	}
	return &v.Database, nil
}
