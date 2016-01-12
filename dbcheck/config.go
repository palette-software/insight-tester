package main

import (
    "os"
    "io/ioutil"
    "encoding/xml"
)

type Server struct {
    Host string `xml:"host,attr"`
    Port int `xml:"port,attr"`
}

type User struct {
    Login string `xml:"login,attr"`
    Password string `xml:"password,attr"`
}

type Table struct {
    Name string `xml:"name,attr"`
}

type Database struct {
    Name string `xml:"name,attr"`
    Type string `xml:"type,attr"`
    Server Server
    User User
    Table Table
}

type Config struct {
    Database Database
}

func getResultDBConfig(fileName string) (Config) {
    var v Config
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
     err = xml.Unmarshal(b, &v)
     if err != nil {
        Error.Println("Error parsing xml ", err)
        os.Exit(1)
    }
    return v
}
