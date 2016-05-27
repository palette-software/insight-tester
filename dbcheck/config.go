package main

import (
	log "github.com/palette-software/insight-tester/common/logging"
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

func getResultDBConfig(fileName string) Database {
	var v Config
	input, err := os.Open(fileName)
	if err != nil {
		log.Error("Error opening file: ", err)
		os.Exit(1)
	}
	defer input.Close()
	b, err := ioutil.ReadAll(input)
	if err != nil {
		log.Error("Error reading file: ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		log.Error("Error parsing yml", err)
		os.Exit(1)
	}
	if v.Database.Host == "" ||
		v.Database.Port == 0 ||
		v.Database.User == "" ||
		v.Database.Database == "" ||
		v.Database.Password == "" {
		log.Error("Config file does not contain database information.")
		os.Exit(1)
	}
	return v.Database
}
