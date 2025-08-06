package logger

import (
	"log"
	"os"
)

type Logger struct {
	logInfo *log.Logger
	logErr  *log.Logger
}

func (l *Logger) Info(v ...interface{}) {
	l.logInfo.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.logErr.Println(v...)
}

func newLogger(logInfo, logErr *log.Logger) *Logger {
	return &Logger{logInfo: logInfo, logErr: logErr}
}

func SetupLogger(path string) *Logger {
	flags := log.LstdFlags | log.Lshortfile

	fileInfo, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileErr, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	logInfo := log.New(fileInfo, "INFO:\t", flags)
	logErr := log.New(fileErr, "ERROR:\t", flags)

	log := newLogger(logInfo, logErr)

	return log
}
