package db_connector

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	DbConnector         DbConnector `yaml:"DbConnector"`
	SplunkServerAddress string      `yaml:"SplunkServerAddress"`
	SplunkToken         string      `yaml:"SplunkToken"`
	SplunkCustomer      string      `yaml:"SplunkCustomer"`
}

func ParseConfig(fileName string) (*Config, error) {
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
	if v.DbConnector.Host == "" ||
		v.DbConnector.Port == 0 ||
		v.DbConnector.User == "" ||
		v.DbConnector.Schema == "" ||
		v.DbConnector.Database == "" ||
		v.DbConnector.Password == "" {
		return nil, fmt.Errorf("Config file does not contain database information. %v", v.DbConnector)
	}
	return &v, nil
}
