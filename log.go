package main

import (
	"io"
	"log"
	"os"
)

// customized log
var (
	Trace   *log.Logger // For anything
	Info    *log.Logger // For common error
	Warning *log.Logger // For notify
	Error   *log.Logger // For critical problem
)

func init() {
	// initialize custom error logging
	traceFile, err := os.OpenFile("logs/trace.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatalln("Failed to open error log file:", err)
	}

	errorFile, err := os.OpenFile("logs/errors.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(io.MultiWriter(traceFile, os.Stderr),
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(errorFile, os.Stderr),
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(errorFile, os.Stderr),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errorFile, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
