package paw

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cnested"
	"github.com/sirupsen/logrus"
)

// nested "github.com/antonfisher/nested-logrus-formatter"

var (
	Logger = &logrus.Logger{
		Out:          os.Stdout, //os.Stderr,
		ReportCaller: true,
		Formatter:    cnestedFMT,
		// Level:        logrus.InfoLevel,
		Level: logrus.WarnLevel,
	}
	// NestedFormatter ...
	cnestedFMT = cnested.DefaultFormat
	// cnestedFMT = &cnested.Formatter{
	// 	HideKeys: false,
	// 	// FieldsOrder:     []string{"component", "category"},
	// 	NoColors:       false,
	// 	NoFieldsColors: false,
	// 	// TimestampFormat: "2006-01-02 15:04:05",
	// 	TimestampFormat: "060102-150405.000",
	// 	TrimMessages:    true,
	// 	CallerFirst:     true,
	// 	CustomCallerFormatter: func(f *runtime.Frame) string {
	// 		s := strings.Split(f.Function, ".")
	// 		funcName := s[len(s)-1]
	// 		name := filepath.Base(f.File)
	// 		cname := NewLSColor(filepath.Ext(name)).Sprint(name)
	// 		cfuncName := color.New(color.FgYellow).Add(color.Bold).Sprint(funcName)
	// 		cln := NewEXAColor("sn").Sprint(f.Line)
	// 		return fmt.Sprintf(" [%s:%s][%s]", cname, cln, cfuncName)
	// 	},
	// }
)

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	// Logger.SetLevel(logrus.InfoLevel)
	// Logger.SetReportCaller(true)
	// Logger.SetOutput(os.Stdout)
	// Logger.SetFormatter(nestedFormatter)
	GologInit(os.Stdout, os.Stdout, os.Stderr, true)
}

// GologInit initializes logger
func GologInit(
	infoHandle io.Writer,
	warnHandle io.Writer,
	errorHandle io.Writer,
	isVerbose bool) {
	if isVerbose {
		Trace = log.New(infoHandle,
			Ctrace.Add(color.Bold).Sprint("[TRACE] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Debug = log.New(infoHandle,
			Cdebug.Add(color.Bold).Sprint("[DEBUG] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(infoHandle,
			Cinfo.Add(color.Bold).Sprint("[INFO] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(warnHandle,
			Cwarn.Add(color.Bold).Sprint("[WARNING] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(errorHandle,
			Cerror.Add(color.Bold).Sprint("[ERROR] "),
			log.Ldate|log.Ltime|log.Lshortfile)
		return
	}

	Trace = log.New(infoHandle,
		color.New(color.FgCyan).Add(color.Bold).Sprint("[TRAC] "),
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(infoHandle,
		color.New(color.FgMagenta).Add(color.Bold).Sprint("[DEBU] "),
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		color.New(color.FgYellow).Add(color.Bold).Sprint("[INFO] "),
		log.Ldate|log.Ltime)

	Warning = log.New(warnHandle,
		color.New(color.FgHiMagenta).Add(color.Bold).Sprint("[WARN] "),
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle,
		color.New(color.FgHiRed).Add(color.Bold).Sprint("[ERRO] "),
		log.Ldate|log.Ltime)

}

// LoggerSetFieldsOrder set `nestedFormatter.FieldsOrder`
func LoggerSetFieldsOrder(fields []string) {
	cnestedFMT.FieldsOrder = fields
	Logger.SetFormatter(cnestedFMT)
}

// var (
// 	// Log is a logger use logrus
// 	Log = logrus.New()
// 	// // Info is log.Info
// 	// Info *log.Logger
// 	// // Warn is log.Warn
// 	// Warn *log.Logger
// 	// // Error is log.Error
// 	// Error *log.Logger
// 	// // // Debug ...
// 	// // Debug *log.Logger
// 	// // // Fatal ...
// 	// // Fatal *log.Logger
// )

// func init() {
// 	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// 	// Info = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
// 	// Warn = log.New(os.Stdout, "[Warn] ", log.Ldate|log.Ltime|log.Lshortfile)
// 	// Error = log.New(os.Stderr, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
// 	// Log as JSON instead of the default ASCII formatter.
// 	// Log.SetFormatter(&log.JSONFormatter{})
// 	// Log.SetFormatter(&log.TextFormatter{})
// 	// Log.SetFormatter(&nested.Formatter{
// 	// 	HideKeys: true,
// 	// 	// TimestampFormat: time.RFC3339,
// 	// 	TimestampFormat: time.Now().Local().Format("0102:150405.000"),
// 	// 	// FieldsOrder:     []string{"level", "func", "file", "msg"},
// 	// })
// 	Log.SetFormatter(new(LogFormatter))
// 	logrus.SetFormatter(new(LogFormatter))
// 	// Output to stdout instead of the default stderr
// 	// Can be any io.Writer, see below for File example
// 	Log.SetOutput(os.Stdout)
// 	logrus.SetOutput(os.Stdout)

// 	// Only log the warning severity or above.
// 	// Log.SetLevel(logrus.WarnLevel)
// 	// Log.SetLevel(log.InfoLevel)
// 	// logrus.SetLevel(log.InfoLevel)
// 	Log.SetReportCaller(true)
// 	logrus.SetReportCaller(true)
// }

// //LogFormatter 日誌自定義格式
// type LogFormatter struct{}

// //Format 格式詳情
// func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	timestamp := time.Now().Local().Format("0102-150405.000")
// 	var (
// 		file   string
// 		len    int
// 		method string
// 	)
// 	if entry.Caller != nil {
// 		file = filepath.Base(entry.Caller.File)
// 		len = entry.Caller.Line
// 		method = filepath.Base(entry.Caller.Func.Name())
// 	}
// 	//fmt.Println(entry.Data)
// 	// msg := fmt.Sprintf("%s [%s:%d][GOID:%d][%s] %s\n", timestamp, file, len, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
// 	msg := fmt.Sprintf("%s [%s:%d][%s][%s] %s\n", timestamp, file, len, method, strings.ToUpper(entry.Level.String()), entry.Message)
// 	return []byte(msg), nil
// }

// func getGID() uint64 {
// 	b := make([]byte, 64)
// 	b = b[:runtime.Stack(b, false)]
// 	b = bytes.TrimPrefix(b, []byte("goroutine "))
// 	b = b[:bytes.IndexByte(b, ' ')]
// 	n, _ := strconv.ParseUint(string(b), 10, 64)
// 	return n
// }

// // type Job struct {
// // 	Command string
// // 	*log.Logger
// // }

// // func NewJob(command string) *Job {
// // 	return &Job{command, log.New(os.Stderr, "Job: ", log.Ldate|log.Ltime|log.Lshortfile)}
// // }
