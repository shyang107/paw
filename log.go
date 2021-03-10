package paw

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cnested"
	"github.com/spf13/cast"

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
		color.New(color.FgCyan).Sprint("[TRAC] "),
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(infoHandle,
		color.New(color.FgMagenta).Add(color.Bold).Sprint("[DEBU] "),
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		color.New(color.FgYellow).Sprint("[INFO] "),
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

// -------------------------------------

type ValuePair struct {
	Field      string
	Value      interface{}
	FieldColor *color.Color
	ValueColor *color.Color
}

func NewValuePair(field string, value interface{}) *ValuePair {
	return &ValuePair{
		Field:      field,
		Value:      value,
		FieldColor: Cnop,
		ValueColor: Cvalue,
	}
}

func (v ValuePair) String() string {
	return MesageFieldAndValueC(
		v.Field,
		v.Value,
		Logger.GetLevel(),
		v.FieldColor,
		v.ValueColor,
	)
}

type ValuePairA []*ValuePair

func NewValuePairA(cap int) ValuePairA {
	if cap < 0 {
		cap = 0
	}
	return make(ValuePairA, 0, cap)
}

func (v *ValuePairA) Add(field string, value interface{}) *ValuePairA {
	(*v) = append((*v), NewValuePair(field, value))
	return v
}

func (v ValuePairA) String() string {
	sb := new(strings.Builder)
	for _, vp := range v {
		sb.WriteString(vp.String())
	}
	return sb.String() ///
}

func MesageFieldAndValueC(field string, value interface{}, level logrus.Level, cf, cv *color.Color) string {
	if cf == nil {
		cf = LogLevelColor(level)
	}

	if cv == nil {
		cv = Cvalue
	}
	msg := "[" + cf.Sprintf("%s: ", field)
	msg += cv.Sprintf("%v", value)
	msg += "]"
	return msg
}

func MesageFieldAndValue(field string, value interface{}, level logrus.Level) string {
	return "[" + field + ": " + cast.ToString(value) + "]"
}
