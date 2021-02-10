package paw

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cnested"

	"github.com/sirupsen/logrus"
	// nested "github.com/antonfisher/nested-logrus-formatter"
)

var (
	sb = new(strings.Builder)
	// Logger is logrus.Logger
	// Logger = logrus.New()
	Logger = &logrus.Logger{
		Out:          os.Stdout, //os.Stderr,
		ReportCaller: true,
		Formatter:    nestedFormatter,
		// Level:        logrus.InfoLevel,
		Level: logrus.WarnLevel,
	}
	// NestedFormatter ...
	nestedFormatter = &cnested.Formatter{
		HideKeys: false,
		// FieldsOrder:     []string{"component", "category"},
		NoColors:       false,
		NoFieldsColors: false,
		// TimestampFormat: "2006-01-02 15:04:05",
		TimestampFormat: "0102-150405.000",
		TrimMessages:    true,
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			name := filepath.Base(f.File)
			cname := NewLSColor(filepath.Ext(name)).Sprint(name)
			cfuncName := color.New(color.FgHiGreen).Add(color.Bold).Sprint(funcName)
			cln := NewEXAColor("sn").Sprint(f.Line)
			return fmt.Sprintf(" [%s:%s][%s]", cname, cln, cfuncName)
		},
	}
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
			color.New(color.FgCyan).Add(color.Bold).Sprint("[TRACE] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Debug = log.New(infoHandle,
			color.New(color.FgMagenta).Add(color.Bold).Sprint("[DEBUG] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(infoHandle,
			color.New(color.FgYellow).Add(color.Bold).Sprint("[INFO] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(warnHandle,
			color.New(color.FgHiMagenta).Add(color.Bold).Sprint("[WARNING] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(errorHandle,
			color.New(color.FgHiRed).Add(color.Bold).Sprint("[ERROR] "),
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

// SetLoggerFieldsOrder set `nestedFormatter.FieldsOrder`
func SetLoggerFieldsOrder(fields []string) {
	nestedFormatter.FieldsOrder = fields
	Logger.SetFormatter(nestedFormatter)
}
