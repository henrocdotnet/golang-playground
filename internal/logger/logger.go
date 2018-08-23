package logger

import "log"

var (
	DebugMode = false
)

func Debug(m string, v ...interface{}) {
	if !DebugMode {
		return
	}

	log.Printf(m, v...)
}

