package main

import (
    "log"
    "io"
    "io/ioutil"
    "os"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func InitLog(
            traceHandle io.Writer,
            infoHandle io.Writer,
            warningHandle io.Writer,
            errorHandle io.Writer)  {
    Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
    Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
    InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    if len(os.Args) < 3 {
        Error.Printf("Usage: %s test_config_json palmon_config_xml\n", os.Args[0])
        os.Exit(1)
    }
    exitCode := 0

    tests := getTests(os.Args[1])
    config := getResultDBConfig(os.Args[2])
    for _, test := range tests {
        if !check(config, test) {
            exitCode = 1
        }
    }
    closeDB()
    os.Exit(exitCode)
}


