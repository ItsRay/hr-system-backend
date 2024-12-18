package common

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.Logger.SetPrefix("[INFO] ")
	l.Logger.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Logger.SetPrefix("[WARN] ")
	l.Logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Logger.SetPrefix("[ERROR] ")
	l.Logger.Println(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Logger.SetPrefix("[FATAL] ")
	l.Logger.Println(v...)
	os.Exit(1)
}

func (l *Logger) TimeTrack(start time.Time, name string) {
	duration := time.Since(start)
	l.Info(name, "took", duration)
}
