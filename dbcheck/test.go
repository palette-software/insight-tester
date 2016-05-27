package main

import (
	"encoding/json"
	"fmt"
	log "github.com/palette-software/insight-tester/common/logging"
	"io/ioutil"
	"os"
	"strconv"
)

type Result struct {
	Operation string
	Count     int
}

type Test struct {
	Description string
	Sql         string
	Result      Result
}

func color(this_color int, input string) string {
	return "\033[" + strconv.Itoa(this_color) + "m" + input + "\033[0m"
}

func print(description string, result bool) {
	output := description + ": "
	if result {
		output += "OK"
		output = color(32, output)
	} else {
		output += "Failed!!!!"
		output = color(31, output)
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

func getTests(fileName string) []Test {
	var v []Test
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
	err = json.Unmarshal(b, &v)
	if err != nil {
		log.Error("Error parsing json ", err)
		os.Exit(1)
	}
	return v
}
