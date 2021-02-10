package cnested

//
// based on github.com/antonfisher/nested-logrus-formatter
//
import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	cl := getColorOfLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := new(strings.Builder)

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if f.ShowFullLevel {
		level = "[" + level + "]"
	} else {
		level = "[" + level[:4] + "]"
	}

	if f.NoColors {
		fmt.Fprint(b, level)
	} else {
		cl.Fprint(b, level)
	}
	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	if f.CallerFirst {
		// f.writeCaller(b, entry)
		fmt.Fprint(b, f.sCaller(entry))
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	// write fields
	var sfields string
	if f.FieldsOrder == nil {
		// f.writeFields(b, entry)
		sfields = f.sFields(entry)
	} else {
		// f.writeOrderedFields(b, entry)
		sfields = f.sOrderedFields(entry)
	}
	if f.NoColors {
		fmt.Fprint(b, sfields)
	} else {
		cl.Fprint(b, sfields)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	// write message
	var mesg string
	if f.TrimMessages {
		mesg = strings.TrimSpace(entry.Message)
	} else {
		mesg = entry.Message
	}
	cl.Fprint(b, mesg)

	if !f.CallerFirst {
		fmt.Fprint(b, f.sCaller(entry))
	}

	b.WriteByte('\n')

	return []byte(b.String()), nil
}

func (f *Formatter) sCaller(entry *logrus.Entry) (s string) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			s = f.CustomCallerFormatter(entry.Caller)
		} else {
			s = fmt.Sprintf(
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
	return s
}

func (f *Formatter) sFields(entry *logrus.Entry) (s string) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)
		for _, field := range fields {
			s += f.sField(entry, field)
		}
	}
	return s
}

func (f *Formatter) sOrderedFields(entry *logrus.Entry) (s string) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			s += f.sField(entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.sField(entry, field)
		}
	}
	return s
}

func (f *Formatter) sField(entry *logrus.Entry, field string) (s string) {

	if f.HideKeys {
		s = fmt.Sprintf("[%v]", entry.Data[field])
	} else {
		s = fmt.Sprintf("[%s:%v]", field, entry.Data[field])
	}

	if !f.NoFieldsSpace {
		s += " "
	}
	return s
}

func getColorOfLevel(level logrus.Level) (c *color.Color) {
	switch level {
	case logrus.TraceLevel:
		c = color.New(color.FgCyan)
	case logrus.DebugLevel:
		c = color.New(color.FgMagenta)
	case logrus.WarnLevel:
		c = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		c = color.New([]color.Attribute{38, 5, 220, 1, 48, 5, 160}...) //
	default:
		c = color.New(color.FgBlue)
	}
	return c
}
