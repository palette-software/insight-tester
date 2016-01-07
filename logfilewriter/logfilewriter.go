package main

import (
    "os"
    "io"
    "log"
    "bufio"
    "sync"
    "io/ioutil"
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

func readFile(fileName string, transfer chan string) (error) {
    file, err := os.Open(fileName)
    defer close(transfer)
    if err != nil {
        Error.Printf("Error opening file: %v\n", err)
        return err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    var line string
    for err == nil {
        line, err = reader.ReadString('\n')
        if err == nil {
            transfer <- line
        }
    }
    return nil
}

func writeFile(path string, fileName string, transfer chan string, wg* sync.WaitGroup) (error) {
    defer wg.Done()
    file, err := os.Create(path + fileName)
    if err != nil {
        Error.Printf("Error opening file: %v", err)
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for line := range transfer {
        writer.WriteString(line)
    }
    writer.Flush()
    return nil
}


func main() {
    InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    if len(os.Args) < 3 {
        Error.Printf("Usage: %s source_path destination_path\n", os.Args[0])
        os.Exit(1)
    }

    files, err := ioutil.ReadDir(os.Args[1])
    if err != nil {
        Error.Println("Error while reading source directory: ", err)
    }

    var wg sync.WaitGroup
    for _, file := range files {
        if !file.IsDir() {
            wg.Add(1)
            Info.Println("Processing file: ", file.Name())
            transfer := make(chan string)
            go writeFile(os.Args[2], file.Name(), transfer, &wg)
            go readFile(file.Name(), transfer)
        }
    }
    wg.Wait()
    Trace.Println("Exiting...")
    os.Exit(0)
}


