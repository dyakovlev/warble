package utils

import (
	"log"
	"os"
)

var (
	Info    func(...interface{})
	Warning func(...interface{})
	Error   func(...interface{})
)

func init() {
	Info = makeLogger(log.New(os.Stdout, "INFO: ", log.Lshortfile|log.Ldate|log.Ltime))
	Warning = makeLogger(log.New(os.Stdout, "WARNING: ", log.Lshortfile|log.Ldate|log.Ltime))
	Error = makeLogger(log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.Ldate|log.Ltime))
}

func makeLogger(l *log.Logger) func(v ...interface{}) {
	return func(v ...interface{}) {
		l.Println(v...)
	}
}
