# Basic stuff.

# .PHONY: directories

all: directories file-summer note-monger

file-summer: directories
	go build -o bin/file-summer cmd/file-summer/file-summer.go

file-summer-watch: directories
	find | entr -s "clear && make file-summer && ./bin/file-summer"

note-monger: directories
	go build -o bin/note-monger cmd/note-monger/main.go

note-monger-watch: directories
	find | entr -s "clear; killall -q note-monger; make note-monger && ./bin/note-monger"

directories:
	mkdir -p bin
