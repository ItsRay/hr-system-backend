package common

import (
	"fmt"
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

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.SetPrefix("[INFO] ")
	message := fmt.Sprintf(format, v...)
	l.Logger.Println(message)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[WARN] ")
	message := fmt.Sprintf(format, v...)
	l.Logger.Println(message)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[ERROR] ")
	message := fmt.Sprintf(format, v...)
	l.Logger.Println(message)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.SetPrefix("[FATAL] ")
	message := fmt.Sprintf(format, v...)
	l.Logger.Println(message)
	os.Exit(1)
}

func (l *Logger) TimeTrack(start time.Time, name string) {
	duration := time.Since(start)
	l.Infof(name, "took", duration)
}
