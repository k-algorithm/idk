package filelogger

import (
	"fmt"
	"os"

	"log"
)

//go:generate mockgen -source=logger.go -destination=logger_mock.go -package=statslog Logger

const (
	defaultStatus    = "INFO"
	defaultTagPrefix = "IDK"

	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Logger abstracts FileLogger.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Log(status string, msg string)
}

// FileLogger is a logger that logs to a file.
type FileLogger struct {
	prefix   string
	fileName string
}

func New(prefix string, fileName string) Logger {
	if prefix == "" {
		prefix = defaultTagPrefix
	}
	fileLogger := newFileLogger(prefix, fileName)
	return fileLogger
}

func newFileLogger(prefix string, fileName string) *FileLogger {
	return &FileLogger{prefix: prefix, fileName: fileName}
}

func (l *FileLogger) Info(msg string) {
	coloredStatus := fmt.Sprint(colorGreen + "[INFO]" + colorReset)
	l.Log(coloredStatus, msg)
}

func (l *FileLogger) Warn(msg string) {
	coloredStatus := fmt.Sprint(colorYellow + "[WARN]" + colorReset)
	l.Log(coloredStatus, msg)
}

func (l *FileLogger) Error(msg string) {
	coloredStatus := fmt.Sprint(colorRed + "[ERROR]" + colorReset)
	l.Log(coloredStatus, msg)
}

func (l *FileLogger) Log(status string, msg string) {
	file, err := os.OpenFile(l.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	if status == "" {
		status = defaultTagPrefix
	}
	prefix := fmt.Sprintf("%s[%s] ", status, l.prefix)
	defer file.Close()
	fileLogger := log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)

	if msg == "" {
		fileLogger.Println("Empty message")
		return
	}

	fileLogger.Print(msg)
}
