package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const MODE_SYNC_WAITGROUP = "waitgroup"
const MODE_SYNC_CHANNEL = "channel"

var (
	channelProcess  = make(chan hashfile)
	channelComplete = make(chan hashfile)
	channelDone     = make(chan int)
	debugMode       = false
	directory       = "."
	workerLimit     = 10
	wg              sync.WaitGroup

	countFilesFound     = 0
	countFilesProcessed = 0
	countFilesPrinted   = 0
	modeSync            = MODE_SYNC_CHANNEL
)

func init() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode.")
	flag.StringVar(&directory, "directory", ".", "Directory to scan.")
	flag.IntVar(&workerLimit, "limit", 10, "Hash calculator worker limit.")
	flag.Parse()
}

func main() {
	debugMessage("BEGIN: main")

	// Test scan directory.
	validateDirectoryOrExit(directory)

	// Setup workers.
	go workerPrintResult(channelComplete)
	for i := 1; i <= workerLimit; i++ {
		go workerProcessFile(channelProcess)
	}

	// Find files and send to workers..
	filepath.Walk(directory, findFilesCallback)
	close(channelProcess)

	// Wait until all files have been found, hashed, and printed.
	if modeSync == MODE_SYNC_WAITGROUP {
		fmt.Println("USING SYNC WAITGROUP")
		wg.Wait()
	} else {
		<-channelDone
	}

	fmt.Print("\nTotals:\n")
	fmt.Printf("    Found: %d\n", countFilesFound)
	fmt.Printf("Processed: %d\n", countFilesProcessed)
	fmt.Printf("  Printed: %d\n", countFilesPrinted)
}

// Path walker callback.
// TODO: Error parameter should be checked here.
func findFilesCallback(path string, info os.FileInfo, e error) error {
	// Skip directories.
	if info.IsDir() {
		return nil
	}

	// Create new hashfile.
	hf := hashfile{path: path}

	if modeSync == MODE_SYNC_WAITGROUP {
		wg.Add(1)
	}

	// Add to process queue.
	countFilesFound += 1
	channelProcess <- hf
	debugMessage("BEGIN: findFilesCallback: found %s", path)

	return nil

}

// Prints the path and hash of a processed file.
func workerPrintResult(c chan (hashfile)) {
	for f := range c {
		fmt.Printf("%s: %s\n", f.path, f.hash)
		countFilesPrinted += 1

		if modeSync == MODE_SYNC_WAITGROUP {
			wg.Done()
		} else if countFilesPrinted == countFilesFound { // TODO: This isn't working, if the printer gets ahead of the checksummer.
			channelDone <- 1
			return
		}
	}
}
// Handles hash generation for a file.
func workerProcessFile(c chan (hashfile)) {
	for f := range c {

		bytes, err := ioutil.ReadFile(f.path)
		if err != nil {
			fmt.Printf("ERROR: Could not read file '%s': %s'\n", f.path, err)
			continue
		}

		hash := md5.New()
		hash.Write(bytes)
		// f.hash = fmt.Sprintf("%x", hex.EncodeToString(hash.Sum(nil)))
		f.hash = hex.EncodeToString(hash.Sum(nil))

		countFilesProcessed += 1
		channelComplete <- f
	}
}

// Ensures the path submitted exists and is a directory.
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
