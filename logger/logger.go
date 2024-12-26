package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) log(color string, level string, message string, fields map[string]interface{}) {
	formattedFields := ""
	for key, value := range fields {
		formattedFields += fmt.Sprintf("%s=%v ", key, value)
	}
	log.Printf("%s[%s]%s %s %s\n", color, level, Reset, message, formattedFields)
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(Green, "INFO", message, fields)
}

func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(Yellow, "WARN", message, fields)
}

func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(Red, "ERROR", message, fields)
}

func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	l.log(Red, "FATAL", message, fields)
	os.Exit(1)
}
