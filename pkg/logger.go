package logger

import (
	"fmt"
	"log"
	"os"
	"path"
)

const defaultErrorLogPath = "./logs/err"

func NewErrorFile(filename string) *log.Logger {
	if err := os.MkdirAll(defaultErrorLogPath, 0777); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	f, err := os.OpenFile(path.Join(defaultErrorLogPath, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	err = f.Truncate(0)
	if err != nil {
		log.Fatalf("Failed to truncate file: %v", err)
	}

	newLog := log.New(f, fmt.Sprintf("[%s] ", filename), log.LstdFlags)
	return newLog
}
