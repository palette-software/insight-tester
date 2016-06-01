package main

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Result struct {
	Operation string `yaml:"Operation"`
	Count     int    `yaml:"Count"`
}

type Test struct {
	Description string `yaml:"Description"`
	Sql         string `yaml:"Sql"`
	Result      Result `yaml:"Result"`
}

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

func print(description string, result bool) {
	output := description + ": "
	if result {
		output += "OK"
		output = green(output)
	} else {
		output += "Failed!!!!"
		output = red(output)
	}
	fmt.Println(output)
}

func checkTest(count int, test Test) bool {
	ret := false
	switch {
	case test.Result.Operation == "<":
		ret = count < test.Result.Count
	case test.Result.Operation == "=":
		ret = count == test.Result.Count
	case test.Result.Operation == ">":
		ret = count > test.Result.Count
	}
	print(test.Description, ret)
	return ret
}

func getTests(fileName string) ([]Test, error) {
	var v []Test
	input, err := os.Open(fileName)
	if err != nil {
		return v, err
	}
	defer input.Close()
	b, err := ioutil.ReadAll(input)
	if err != nil {
		return v, err
	}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		return v, err
	}
	return v, nil
}
