package main

import "fmt"
import "os"
import "path/filepath"
import "strings"
import "io"

func MoveFilesButKeepFolders(srcFolder string, dstFolder string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		dstPath := strings.Replace(path, srcFolder, dstFolder, 1)
		if f.IsDir() {
			//fmt.Printf("Visited: %s\n", dstPath)
			os.MkdirAll(dstPath, os.ModePerm)
		} else {
			os.Rename(path, dstPath)
			fileCounter = fileCounter + 1
		}
		return nil
	}
}

func CopyFilesAndFolders(srcFolder string, dstFolder string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {

		dstPath := strings.Replace(path, srcFolder, dstFolder, 1)

		if f.IsDir() {
			os.MkdirAll(dstPath, os.ModePerm)
		} else {
			r, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer r.Close()

			w, err := os.Create(dstPath)
			if err != nil {
				panic(err)
			}

			defer w.Close()

			n, err := io.Copy(w, r)
			fileCounter = fileCounter + 1
			_ = n
			if err != nil {
				panic(err)
			}
		}

		return nil;
	}
}

var fileCounter int


func main() {

    if len(os.Args) < 5 {
        fmt.Println("Usage: ", os.Args[0], " <tmpFolder> <srcFolder> <trgUploadsFolder> <trgForkedFolder>")
        fmt.Println(`
                    tmpFolder: Technical folder for keeping files being processed.
                    srcFolder: The folder that contains the original csv files.
                    trgUploadsFolder: 1st folder for the fork.
                    trgForkedFolder: 2nd folder for the fork.
                  `)
        os.Exit(1)
    }
               
       
    var tmpFolder = os.Args[1]
	var srcFolder = os.Args[2]
	var trgUploadsFolder = os.Args[3]
	var trgForkedFolder = os.Args[4]
    
	fmt.Printf("Job started.\n")
	fileCounter = 0
	filepath.Walk(srcFolder, MoveFilesButKeepFolders(srcFolder, tmpFolder))
	fmt.Printf("Number of files moved from %s to %s: %d \n", srcFolder, tmpFolder, fileCounter)
	var fileCounterMoved int = fileCounter

	fileCounter = 0
	filepath.Walk(tmpFolder, CopyFilesAndFolders(tmpFolder, trgForkedFolder))
	fmt.Printf("Number of files copied from %s to %s: %d \n", tmpFolder, trgForkedFolder, fileCounter)
	var fileCounterForkedCopied int = fileCounter

	fileCounter = 0
	filepath.Walk(tmpFolder, CopyFilesAndFolders(tmpFolder, trgUploadsFolder))
	fmt.Printf("Number of files copied from %s to %s: %d \n", tmpFolder, trgUploadsFolder, fileCounter)
	var fileCounterUploadsCopied int = fileCounter

	os.RemoveAll(tmpFolder)
	fmt.Printf("Tmp folder has been deleted.\n")

	if fileCounterMoved != fileCounterForkedCopied ||
		fileCounterForkedCopied != fileCounterUploadsCopied {
		panic("Number of files in tmp not equal to number of copied files")
	}

	fmt.Printf("Job successfully ended.\n")
}
