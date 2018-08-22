package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"flag"
	"encoding/hex"
)

var (
	channelProcess  = make(chan hashfile)
	channelComplete = make(chan hashfile)
	channelDone     = make(chan int)
	debugMode       = false
	directory       = "."
	workerLimit     = 10
)

func init() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode.")
	flag.StringVar(&directory, "directory", ".", "Directory to scan.")
	flag.IntVar(&workerLimit, "limit", 10, "Hash calculator worker limit.")
	flag.Parse()
}

func main() {
	fmt.Printf("WTF? %s\n", directory)
	debugMessage("BEGIN: main")

	// Test scan directory.
	validateDirectoryOrExit(directory)

	// Setup workers.
	go workerPrintResult(channelComplete)
	for i := 1; i <= workerLimit; i++ {
		go workerProcessFile(channelProcess);
	}

	// Find files and send to workers..
	filepath.Walk(directory, findFilesCallback)

	for {
		<- channelDone
	}

	/*
	close(channelProcess)
	close(channelComplete)
	*/
}

func findFilesCallback(path string, info os.FileInfo, e error) error {
	// Skip directories.
	if info.IsDir() {
		return nil
	}

	// Create new hashfile.
	hf := hashfile{path: path}

	// Add to process queue.
	channelProcess <- hf
	debugMessage("BEGIN: findFilesCallback: found %s", path)

	return nil

}

func workerPrintResult(c chan (hashfile)) {
	debugMessage("BEGIN: workerPrintResult")
	for {
		f := <- channelComplete
		fmt.Printf("%s: %s\n", f.path, f.hash)
	}
}

func workerProcessFile(c chan (hashfile)) {
	debugMessage("BEGIN: workerProcessFile")
	for {
		f := <-c

		bytes, err := ioutil.ReadFile(f.path)
		if err != nil {
			fmt.Printf("ERROR: Could not read file '%s': %s'\n", f.path, err)
			continue
		}

		hash := md5.New()
		hash.Write(bytes)
		// f.hash = fmt.Sprintf("%x", hex.EncodeToString(hash.Sum(nil)))
		f.hash = hex.EncodeToString(hash.Sum(nil))

		channelComplete <- f
	}
}

func validateDirectoryOrExit(p string) {
	v, err := os.Stat(p)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("ERROR: directory '%s' not found.\n", p)
		} else {
			fmt.Printf("ERROR: error accessing directory '%s': %s'\n", p, err)
		}
		os.Exit(1)
	}

	if !v.IsDir() {
		fmt.Printf("ERROR: '%s' is not a directory\n", p)
		os.Exit(1)
	}
}

func debugMessage(m string, v ...interface{}) {
	if !debugMode {
		return
	}

	log.Printf(m, v...)
}

// Types.
// ------

type hashfile struct {
	path string
	hash string
}
