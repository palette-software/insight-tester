package main

import (
    "os"
    "io/ioutil"
    "strconv"
    "encoding/json"
)

type Result struct {
    Operation string 
    Count int
}

type Test struct {
    Description string
    Sql string
    Result Result
}


func color(this_color int, input string) (string) {
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
    Info.Println(output)
}


func checkTest(count int, test Test) (bool) {
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

func getTests(fileName string) ([]Test) {
    var v []Test
    input, err := os.Open(fileName)
    if err != nil {
        Error.Println("Error opening file: ", err)
        os.Exit(1)
    }
    defer input.Close()
    b, err := ioutil.ReadAll(input)
    if err != nil {
        Error.Println("Error reading file: ", err)
        os.Exit(1)
    }
     err = json.Unmarshal(b, &v)
     if err != nil {
        Error.Println("Error parsing xml ", err)
        os.Exit(1)
    }
    return v
}

