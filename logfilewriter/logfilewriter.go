package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"sync"
)

import log "github.com/palette-software/insight-tester/common/logging"

func readFile(fileName string, transfer chan string) error {
	log.Infof("Reading file: %v", fileName)
	file, err := os.Open(fileName)
	defer close(transfer)
	if err != nil {
		log.Errorf("Error opening file: %v\n", err)
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

func writeFile(fileName string, transfer chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("Error opening file: %v", err)
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
	log.AddTarget(os.Stdout, log.LevelDebug)
	if len(os.Args) < 3 {
		log.Errorf("Usage: %s source_path destination_path\n", os.Args[0])
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Error("Error while reading source directory: ", err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			log.Info("Processing file: ", file.Name())
			transfer := make(chan string)
			go writeFile(os.Args[2]+"/"+file.Name(), transfer, &wg)
			go readFile(os.Args[1]+"/"+file.Name(), transfer)
		}
	}
	wg.Wait()
	log.Info("Exiting...")

	os.Exit(0)
}
