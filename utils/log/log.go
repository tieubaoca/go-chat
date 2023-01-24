package log

import (
	"log"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func New(file *os.File) {
	InfoLogger.SetOutput(file)
	WarnLogger.SetOutput(file)
	ErrorLogger.SetOutput(file)
	FatalLogger.SetOutput(file)
}
