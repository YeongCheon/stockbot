package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	log.Logger
}

func (l *Logger) Trace(message string) {
	l.Println(message)
}

func NewLogger() *Logger {
	l := &Logger{}
	logFile, err := os.Create("stockbot.log")
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	l.SetOutput(multiWriter)
	l.SetPrefix("INFO: ")
	l.SetFlags(log.LstdFlags)

	return l
}
