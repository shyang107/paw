package paw

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/cnested"

	"github.com/sirupsen/logrus"
)

// nested "github.com/antonfisher/nested-logrus-formatter"

var (
	Logger = &logrus.Logger{
		Out:          os.Stdout, //os.Stderr,
		ReportCaller: true,
		Formatter:    CnestedFMT,
		// Level:        logrus.InfoLevel,
		Level: logrus.WarnLevel,
	}
	// NestedFormatter ...
	CnestedFMT = cnested.DefaultFormat
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
	tracePrefix  = "[TRAC] "
	debugPrefix  = "[DEBU] "
	infoPrefix   = "[INFO] "
	warnPrefix   = "[WARN] "
	errorPrefix  = "[ERRO] "
	fatalPrefix  = "[FATA] "
	tracePrefixF = "[TRACE] "
	debugPrefixF = "[DEBUG] "
	infoPrefixF  = "[INFO] "
	warnPrefixF  = "[WARNING] "
	errorPrefixF = "[ERROR] "
	fatalPrefixF = "[FATAL] "

	tracePrefixC  = Ctrace.Sprint(tracePrefix) + " "
	debugPrefixC  = Cdebug.Sprint(debugPrefix) + " "
	infoPrefixC   = Cinfo.Sprint(infoPrefix) + " "
	warnPrefixC   = Cwarn.Sprint(warnPrefix) + " "
	errorPrefixC  = Cerror.Sprint(errorPrefix) + " "
	fatalPrefixC  = Cerror.Sprint(fatalPrefix) + " "
	tracePrefixFC = Ctrace.Sprint(tracePrefixF) + " "
	debugPrefixFC = Cdebug.Sprint(debugPrefixF) + " "
	infoPrefixFC  = Cinfo.Sprint(infoPrefixF) + " "
	warnPrefixFC  = Cwarn.Sprint(warnPrefixF) + " "
	errorPrefixFC = Cerror.Sprint(errorPrefixF) + " "
	fatalPrefixFC = Cerror.Sprint(fatalPrefixF) + " "

	traceLogo = cnested.Logos[logrus.TraceLevel] + " "
	debugLogo = cnested.Logos[logrus.DebugLevel] + " "
	infoLogo  = cnested.Logos[logrus.InfoLevel] + " "
	warnLogo  = cnested.Logos[logrus.WarnLevel] + " "
	errorLogo = cnested.Logos[logrus.ErrorLevel] + " "
	fatalLogo = cnested.Logos[logrus.FatalLevel] + " "
)

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger
)

func init() {
	// Logger.SetLevel(logrus.InfoLevel)
	// Logger.SetReportCaller(true)
	// Logger.SetOutput(os.Stdout)
	// Logger.SetFormatter(nestedFormatter)
	CnestedFMT.IsLogo = true
	GologInit(os.Stderr, os.Stderr, os.Stderr, false)
}

// GologInit initializes logger
func GologInit(
	infoHandle io.Writer,
	warnHandle io.Writer,
	errorHandle io.Writer,
	isVerbose bool) {
	// cnestedFMT.IsLogo = true
	// if cnestedFMT.IsLogo {
	// 	tracePrefixC = traceLogo + tracePrefixC
	// 	debugPrefixC = debugLogo + debugPrefixC
	// 	infoPrefixC = infoLogo + infoPrefixC
	// 	warnPrefixC = warnLogo + warnPrefixC
	// 	errorPrefixC = errorLogo + errorPrefixC
	// 	tracePrefixFC = traceLogo + tracePrefixFC
	// 	debugPrefixFC = debugLogo + debugPrefixFC
	// 	infoPrefixFC = infoLogo + infoPrefixFC
	// 	warnPrefixFC = warnLogo + warnPrefixFC
	// 	errorPrefixFC = errorLogo + errorPrefixFC
	// }
	if isVerbose {
		Trace = log.New(infoHandle, tracePrefixFC,
			log.Ldate|log.Ltime|log.Lshortfile)

		Debug = log.New(infoHandle, debugPrefixFC,
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(infoHandle, infoPrefixFC,
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(warnHandle, warnPrefixFC,
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(errorHandle, errorPrefixFC,
			log.Ldate|log.Ltime|log.Lshortfile)

		Fatal = log.New(errorHandle, fatalPrefixFC,
			log.Ldate|log.Ltime)
		return
	}

	Trace = log.New(infoHandle, tracePrefixC,
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(infoHandle, debugPrefixC,
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle, infoPrefixC,
		log.Ldate|log.Ltime)

	Warning = log.New(warnHandle, warnPrefixC,
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle, errorPrefixC,
		log.Ldate|log.Ltime)

	Fatal = log.New(errorHandle, fatalPrefixC,
		log.Ldate|log.Ltime)

}

// LoggerSetFieldsOrder set `nestedFormatter.FieldsOrder`
func LoggerSetFieldsOrder(fields []string) {
	CnestedFMT.FieldsOrder = fields
	Logger.SetFormatter(CnestedFMT)
}

// -------------------------------------

type ValuePair struct {
	Field      string
	Value      interface{}
	LeveL      logrus.Level
	FieldColor *color.Color
	ValueColor *color.Color
}

func NewValuePair(field string, value interface{}) *ValuePair {
	return &ValuePair{
		Field:      field,
		Value:      value,
		LeveL:      Logger.GetLevel(),
		FieldColor: nil,
		ValueColor: Cvalue,
	}
}
func NewValuePairWith(field string, value interface{}, level logrus.Level) *ValuePair {
	return &ValuePair{
		Field:      field,
		Value:      value,
		LeveL:      level,
		FieldColor: nil,
		ValueColor: Cvalue,
	}
}

func (v ValuePair) String() string {
	if v.FieldColor == nil {
		v.FieldColor = LogLevelColor(v.LeveL)
	}
	if v.ValueColor == nil {
		v.ValueColor = Cvalue
	}
	msg := "[" + v.FieldColor.Sprintf("%s:", v.Field)
	msg += v.ValueColor.Sprintf("%v", v.Value)
	msg += "]"
	return msg
}

func (v *ValuePair) SetLevel(level logrus.Level) *ValuePair {
	v.LeveL = level
	return v
}

func (v *ValuePair) SetFieldLevelColor(level logrus.Level) *ValuePair {
	v.FieldColor = LogLevelColor(level)
	return v
}

func (v *ValuePair) SetFieldColor(c *color.Color) *ValuePair {
	v.FieldColor = c
	return v
}

func (v *ValuePair) SetValueColor(c *color.Color) *ValuePair {
	v.FieldColor = c.Add([]color.Attribute{4}...)
	return v
}

type ValuePairA []*ValuePair

func NewValuePairA(cap int) ValuePairA {
	if cap < 0 {
		cap = 0
	}
	return make(ValuePairA, 0, cap)
}

func (v ValuePairA) String() string {
	sb := new(strings.Builder)
	for _, vp := range v {
		sb.WriteString(vp.String())
	}
	return sb.String() ///
}

func (v ValuePairA) Add(field string, value interface{}) ValuePairA {
	v = append(v, NewValuePair(field, value))
	return v
}
func (v ValuePairA) AddV(p *ValuePair) ValuePairA {
	return v.AddA(p)
}

func (v ValuePairA) AddA(vps ...*ValuePair) ValuePairA {
	if len(vps) == 0 {
		return v
	}
	for _, vp := range vps {
		v = append(v, vp)
	}
	return v
}

func (v ValuePairA) SetLevel(level logrus.Level) ValuePairA {
	for _, vp := range v {
		vp.LeveL = level
	}
	return v
}

func (v ValuePairA) ToLogrusFields() logrus.Fields {
	var fds logrus.Fields
	for _, p := range v {
		fds[p.Field] = p.Value
	}
	return fds
}

func (v ValuePairA) SprintSep(sep string) string {
	b := new(strings.Builder)
	nv := len(v)
	for i := 0; i < nv-1; i++ {
		b.WriteString(v[i].String() + sep)
	}
	b.WriteString(v[nv-1].String())
	return b.String()
}

func (v ValuePairA) StringFunc(fc func(p *ValuePair) string) string {
	b := new(strings.Builder)
	for _, p := range v {
		b.WriteString(fc(p))
	}
	return b.String()
}

func MesageFieldAndValueC(field string, value interface{}, level logrus.Level, cf, cv *color.Color) string {
	if cf == nil {
		cf = LogLevelColor(level)
		// if level == logrus.InfoLevel {
		// 	cf = COdd
		// } else {
		// 	cf = LogLevelColor(level)
		// }
	}

	if cv == nil {
		cv = Cvalue
	}
	msg := "[" + cf.Sprintf("%s:", field)
	msg += cv.Sprintf("%v", value)
	msg += "]"
	return msg
}

func MesageFieldAndValue(field string, value interface{}, level logrus.Level) string {
	return "[" + field + ": " + cast.ToString(value) + "]"
}
