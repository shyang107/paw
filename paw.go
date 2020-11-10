package paw

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	// nested "github.com/antonfisher/nested-logrus-formatter"
	nested "github.com/antonfisher/nested-logrus-formatter"
)

var (
	sb = strings.Builder{}
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
		HideKeys: true,
		// FieldsOrder:     []string{"component", "category"},
		// TimestampFormat: "2006-01-02 15:04:05",
		TimestampFormat: "0102-150405.000",
		TrimMessages:    true,
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", filepath.Base(f.File), f.Line, funcName)
		},
	}
)

func init() {
	// Logger.SetLevel(logrus.InfoLevel)
	// Logger.SetReportCaller(true)
	// Logger.SetOutput(os.Stdout)
	// Logger.SetFormatter(nestedFormatter)
}

// SetLoggerFieldsOrder set `nestedFormatter.FieldsOrder`
func SetLoggerFieldsOrder(fields []string) {
	nestedFormatter.FieldsOrder = fields
	Logger.SetFormatter(nestedFormatter)
}
