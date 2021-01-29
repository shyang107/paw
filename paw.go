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

	"github.com/sirupsen/logrus"
	// nested "github.com/antonfisher/nested-logrus-formatter"
	nested "github.com/antonfisher/nested-logrus-formatter"
)

var (
	sb = new(strings.Builder)
	// Logger is logrus.Logger
	// Logger = logrus.New()
	Logger = &logrus.Logger{
		Out:          os.Stderr,
		ReportCaller: true,
		Formatter:    nestedFormatter,
		Level:        logrus.InfoLevel,
	}
	// NestedFormatter ...
	nestedFormatter = &nested.Formatter{
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
			cp := NewLSColor(filepath.Ext(name))
			cname := cp.Sprint(name)
			cfuncName := color.New(color.FgHiGreen).Sprint(funcName)
			cln := NewEXAColor("sn").Sprint(f.Line)
			return fmt.Sprintf(" [%s:%s][%s]", cname, cln, cfuncName)
		},
	}
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
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
		Info = log.New(infoHandle,
			color.New(color.FgBlue).Add(color.Bold).Sprint("[INFO] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(warnHandle,
			color.New(color.FgYellow).Add(color.Bold).Sprint("[WARNING] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(errorHandle,
			color.New(color.FgRed).Add(color.Bold).Sprint("[ERROR] "),
			log.Ldate|log.Ltime|log.Lshortfile)

		Debug = log.New(infoHandle,
			color.New(color.FgHiMagenta).Add(color.Bold).Sprint("[DEBUG] "),
			log.Ldate|log.Ltime|log.Lshortfile)
		return
	}
	Info = log.New(infoHandle,
		color.New(color.FgBlue).Add(color.Bold).Sprint("[INFO] "),
		log.Ldate|log.Ltime)

	Warning = log.New(warnHandle,
		color.New(color.FgYellow).Add(color.Bold).Sprint("[WARN] "),
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle,
		color.New(color.FgRed).Add(color.Bold).Sprint("[ERRO] "),
		log.Ldate|log.Ltime)

	Debug = log.New(infoHandle,
		color.New(color.FgHiMagenta).Add(color.Bold).Sprint("[DEBU] "),
		log.Ldate|log.Ltime|log.Lshortfile)
}

// SetLoggerFieldsOrder set `nestedFormatter.FieldsOrder`
func SetLoggerFieldsOrder(fields []string) {
	nestedFormatter.FieldsOrder = fields
	Logger.SetFormatter(nestedFormatter)
}
