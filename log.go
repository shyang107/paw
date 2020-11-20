package paw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Log is a logger use logrus
	Log = logrus.New()
	// // Info is log.Info
	// Info *log.Logger
	// // Warn is log.Warn
	// Warn *log.Logger
	// // Error is log.Error
	// Error *log.Logger
	// // // Debug ...
	// // Debug *log.Logger
	// // // Fatal ...
	// // Fatal *log.Logger
)

func init() {
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Info = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Warn = log.New(os.Stdout, "[Warn] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Error = log.New(os.Stderr, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Log as JSON instead of the default ASCII formatter.
	// Log.SetFormatter(&log.JSONFormatter{})
	// Log.SetFormatter(&log.TextFormatter{})
	// Log.SetFormatter(&nested.Formatter{
	// 	HideKeys: true,
	// 	// TimestampFormat: time.RFC3339,
	// 	TimestampFormat: time.Now().Local().Format("0102:150405.000"),
	// 	// FieldsOrder:     []string{"level", "func", "file", "msg"},
	// })
	Log.SetFormatter(new(LogFormatter))
	logrus.SetFormatter(new(LogFormatter))
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Log.SetOutput(os.Stdout)
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// Log.SetLevel(logrus.WarnLevel)
	// Log.SetLevel(log.InfoLevel)
	// logrus.SetLevel(log.InfoLevel)
	Log.SetReportCaller(true)
	logrus.SetReportCaller(true)
}

//LogFormatter 日誌自定義格式
type LogFormatter struct{}

//Format 格式詳情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("0102-150405.000")
	var (
		file   string
		len    int
		method string
	)
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
		method = filepath.Base(entry.Caller.Func.Name())
	}
	//fmt.Println(entry.Data)
	// msg := fmt.Sprintf("%s [%s:%d][GOID:%d][%s] %s\n", timestamp, file, len, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
	msg := fmt.Sprintf("%s [%s:%d][%s][%s] %s\n", timestamp, file, len, method, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// type Job struct {
// 	Command string
// 	*log.Logger
// }

// func NewJob(command string) *Job {
// 	return &Job{command, log.New(os.Stderr, "Job: ", log.Ldate|log.Ltime|log.Lshortfile)}
// }
