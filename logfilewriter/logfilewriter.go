package main

import (
    "os"
    "bufio"
    "sync"
    "io/ioutil"
)

import log "github/insight-tester/common/logging"

func readFile(fileName string, transfer chan string) (error) {
    log.Info.Printf("Reading file: %v", fileName)
    file, err := os.Open(fileName)
    defer close(transfer)
    if err != nil {
        log.Error.Printf("Error opening file: %v\n", err)
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

func writeFile(fileName string, transfer chan string, wg* sync.WaitGroup) (error) {
    defer wg.Done()
    file, err := os.Create(fileName)
    if err != nil {
        log.Error.Printf("Error opening file: %v", err)
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
    log.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
    if len(os.Args) < 3 {
        log.Error.Printf("Usage: %s source_path destination_path\n", os.Args[0])
        os.Exit(1)
    }

    files, err := ioutil.ReadDir(os.Args[1])
    if err != nil {
        log.Error.Println("Error while reading source directory: ", err)
    }

    var wg sync.WaitGroup
    for _, file := range files {
        if !file.IsDir() {
            wg.Add(1)
            log.Info.Println("Processing file: ", file.Name())
            transfer := make(chan string)
            go writeFile(os.Args[2] + "/" + file.Name(), transfer, &wg)
            go readFile(os.Args[1] + "/" + file.Name(), transfer)
        }
    }
    wg.Wait()
    log.Info.Println("Exiting...")

    os.Exit(0)
}


