package paw

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// TableFormat define the format used to print out
type TableFormat struct {
	Fields          []string
	LenFields       []int
	Aligns          []Align
	Padding         string
	Sep             string
	TopChar         string
	topBanner       string
	MiddleChar      string
	midBanner       string
	BottomChar      string
	botBanner       string
	beforeMsg       string
	afterMsg        string
	writer          io.Writer
	isAbbrSymbol    bool
	isPrepare       bool
	isPrepareBefore bool
	isPrepareAfter  bool
}

// Align is id that indicate alignment of head-column
type Align int

const (
	space      = " "
	abbrSymbol = "»"
	// AlignLeft align left
	AlignLeft Align = iota
	// AlignCenter align center
	AlignCenter
	// AlignRight align right
	AlignRight
)

// NewTableFormat return a instance of TableFormat
func NewTableFormat() *TableFormat {
	return &TableFormat{
		Fields:          []string{},
		LenFields:       []int{},
		Aligns:          []Align{},
		Padding:         "",
		Sep:             space,
		TopChar:         "=",
		MiddleChar:      "-",
		BottomChar:      "=",
		writer:          os.Stdout,
		isAbbrSymbol:    false,
		isPrepare:       false,
		isPrepareBefore: false,
		isPrepareAfter:  false,
	}
}

var (
// isAbbrSymbol    bool
// isPrepare       bool
// isPrepareBefore bool
// isPrepareAfter  bool
)

// SetBeforeMessage set message to show before table
func (t *TableFormat) SetBeforeMessage(msg string) {
	t.beforeMsg = msg
	if len(t.Padding) > 0 {
		t.beforeMsg = t.Padding + t.beforeMsg
	}
	t.isPrepareBefore = true
}

// SetAfterMessage set message to show after table
func (t *TableFormat) SetAfterMessage(msg string) {
	t.afterMsg = msg
	if len(t.Padding) > 0 {
		t.afterMsg = t.Padding + t.afterMsg
	}
	t.isPrepareAfter = true
}

// Prepare initialize `TableFormat`
func (t *TableFormat) Prepare(w io.Writer) {
	t.check()
	t.writer = w
	t.isPrepare = true
}

func (t *TableFormat) check() {
	if len(t.Sep) == 0 {
		t.Sep = space
	}
	t.checkFields()
	t.checkAlign()
	t.setBanner()
}

func (t *TableFormat) checkAlign() {
	la := len(t.Aligns)
	lf := len(t.Fields)
	switch {
	case la < 1:
		t.Aligns = make([]Align, lf)
		for i := 0; i < len(t.Fields); i++ {
			t.Aligns[i] = AlignRight
		}
	case la > lf:
		t.Aligns = t.Aligns[:lf]
	case la < lf && la > 1:
		tmp := make([]Align, lf)
		copy(tmp, t.Aligns)
		for i := la; i < lf; i++ {
			tmp[i] = AlignRight
		}
		t.Aligns = tmp
	}
}

func (t *TableFormat) checkBannerChar() {
	if len(t.TopChar) != 1 {
		t.TopChar = "="
	}
	if len(t.MiddleChar) != 1 {
		t.MiddleChar = "-"
	}
	if len(t.BottomChar) != 1 {
		t.BottomChar = "="
	}
	if len(t.Sep) == 0 {
		t.Sep = space
	}
}

func (t *TableFormat) setBanner() {
	t.checkBannerChar()
	llf := len(t.LenFields)
	lsep := len(t.Sep)
	tlen := 0
	for i := 0; i < llf-1; i++ {
		tlen += t.LenFields[i]
		tlen += lsep
	}
	tlen += t.LenFields[llf-1]
	t.topBanner = strings.Repeat(t.TopChar, tlen)
	t.botBanner = strings.Repeat(t.BottomChar, tlen)
	sb.Reset()
	for i := 0; i < llf-1; i++ {
		sb.WriteString(strings.Repeat(t.MiddleChar, t.LenFields[i]))
		sb.WriteString(t.Sep)
	}
	sb.WriteString(strings.Repeat(t.MiddleChar, t.LenFields[llf-1]))
	t.midBanner = sb.String()
	sb.Reset()
	if len(t.Padding) > 0 {
		t.topBanner = t.Padding + t.topBanner
		t.midBanner = t.Padding + t.midBanner
		t.botBanner = t.Padding + t.botBanner
	}
}

func (t *TableFormat) checkFields() {
	lf := len(t.Fields)
	llf := len(t.LenFields)
	switch {
	case llf < lf:
		for i := llf; i < lf; i++ {
			hc, ac := CountPlaceHolder(t.Fields[i])
			t.LenFields = append(t.LenFields[:llf:llf], hc+ac)
		}
	case llf > lf:
		t.LenFields = t.LenFields[:lf]
	}
	for i := 0; i < lf; i++ {
		if t.LenFields[i] < len(t.Fields[i]) {
			t.LenFields[i] = len(t.Fields[i])
		}
	}
}

func (t *TableFormat) getRowString(fields []string, widths []int, aligns []Align, sep string, padding string) string {
	sb.Reset()
	lf := len(fields)
	for i := 0; i < lf; i++ {
		v := fields[i]
		wd := widths[i]
		v = GetAbbrString(v, wd, "»")
		nh, na := CountPlaceHolder(v)
		al := aligns[i]
		s := ""
		switch al {
		case AlignLeft:
			s = v + strings.Repeat(space, wd-nh-na)
		case AlignRight:
			s = strings.Repeat(space, wd-nh-na) + v
		case AlignCenter:
			lv := nh + na
			nr := (wd - lv) / 2
			nl := wd - lv - nr
			s = strings.Repeat(space, nl) + v + strings.Repeat(space, nr)
		}
		sb.WriteString(s + sep)
	}
	str := sb.String()
	if !t.isAbbrSymbol {
		t.isAbbrSymbol = strings.Contains(str, abbrSymbol)
	}
	return padding + str
}

// PrintSart print out head-section in `t.Writer`
func (t *TableFormat) PrintSart() error {
	if !t.isPrepare {
		t.Prepare(os.Stdout)
	}
	if t.isPrepareBefore {
		fmt.Fprintln(t.writer, t.beforeMsg)
	}
	fmt.Fprintln(t.writer, t.topBanner)
	fmt.Fprintln(t.writer, t.getRowString(t.Fields, t.LenFields, t.Aligns, " ", t.Padding))
	fmt.Fprintln(t.writer, t.midBanner)
	return nil
}

// PrintRow print row into `t.writer`
func (t *TableFormat) PrintRow(rows ...interface{}) {
	var srows []string
	for _, r := range rows {
		srows = append(srows, fmt.Sprint(r))
	}
	fmt.Fprintln(t.writer, t.getRowString(srows, t.LenFields, t.Aligns, t.Sep, t.Padding))
}

// PrintEnd print end-section into `t.writer`
func (t *TableFormat) PrintEnd() {
	if t.isAbbrSymbol {
		fmt.Fprintln(t.writer, strings.ReplaceAll(t.botBanner, t.BottomChar, t.MiddleChar))
		fmt.Fprintln(t.writer, t.Padding+"* '"+abbrSymbol+"' : abbreviation of term")
	}
	fmt.Fprintln(t.writer, t.botBanner)
	if t.isPrepareAfter {
		fmt.Fprintln(t.writer, t.afterMsg)
	}
}
