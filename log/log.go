package log

import (
	"log"
	"os"
)

var (
	// Info is log.Info
	Info *log.Logger
	// Warn is log.Warn
	Warn *log.Logger
	// Error is log.Error
	Error *log.Logger
	// // Debug ...
	// Debug *log.Logger
	// // Fatal ...
	// Fatal *log.Logger
)

func init() {
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	Info = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stdout, "[Warn] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// type Job struct {
// 	Command string
// 	*log.Logger
// }

// func NewJob(command string) *Job {
// 	return &Job{command, log.New(os.Stderr, "Job: ", log.Ldate|log.Ltime|log.Lshortfile)}
// }
