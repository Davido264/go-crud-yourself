package logger

import (
	"log"
)

var wsLogger *log.Logger = nil

type ManagerLogger struct {
	Logch chan<- []byte
}

func (m ManagerLogger) Write(p []byte) (int, error) {
	m.Logch <- p
	return len(p), nil
}

func LocalOnlyPrintf(format string, v ...any) {
	log.Printf(format, v...)
}

func Printf(format string, v ...any) {
	if wsLogger != nil {
		wsLogger.SetPrefix("INFO")
		wsLogger.Printf(format, v...)
		wsLogger.SetPrefix("")
	}
	log.Printf(format, v...)
}

func Fatalf(format string, v ...any) {
	if wsLogger != nil {
		wsLogger.SetPrefix("ERROR")
		wsLogger.Printf(format, v...)
		wsLogger.SetPrefix("")
	}
	log.Fatalf(format, v...)
}

func Panic(err error) {
	if wsLogger != nil {
		wsLogger.SetPrefix("PANIC")
		wsLogger.Println(err)
		wsLogger.SetPrefix("")
	}
	log.Panic(err)
}

func RegisterLogger(logch chan<- []byte) {
	if wsLogger == nil {
		wsLogger = log.New(&ManagerLogger{Logch: logch}, "", log.LstdFlags)
	}
}
