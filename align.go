package paw

import (
	"github.com/shyang107/paw/cast"
)

// Align - Custom type to hold value for alignment option from 1-3
type Align int

// Declare related constants for each alignment starting with index 1
const (
	// AlignLeft align left
	AlignLeft Align = iota + 1
	// AlignCenter align center
	AlignCenter
	// AlignRight align right
	AlignRight
)

func (a Align) String() string {
	return [...]string{"Left", "Center", "Right"}[a-1]
	// switch a {
	// case AlignLeft:
	// 	return "left"
	// case AlignCenter:
	// 	return "center"
	// case AlignRight:
	// 	return "right"
	// default:
	// 	return "Unknown"
	// }
}

func (a Align) EnumIndex() int {
	return int(a)
}

// ToString return the string of value with filling spaces according to Align
func (a Align) ToString(value interface{}, width int) string {
	v := cast.ToString(value)
	switch a {
	case AlignLeft:
		return FillRight(v, width)
	case AlignRight:
		return FillLeft(v, width)
	default: //AlignCenter
		return FillLeftRight(v, width)
	}
}

// ToString return the string of color string cvalue with filling spaces according to Align
func (a Align) ToStringC(cvalue interface{}, width int) string {
	s := cast.ToString(cvalue)
	v := StripANSI(s)
	switch a {
	case AlignLeft:
		sp := lrspace(v, width)
		return s + sp
	case AlignRight:
		sp := lrspace(v, width)
		return sp + s
	default: //AlignCenter
		lsp, rsp := cspace(v, width)
		return lsp + s + rsp
	}
}

func lrspace(s string, w int) string {
	width := StringWidth(s)
	count := w - width
	if count > 0 {
		return Spaces(count)
	}
	return ""
}

func cspace(s string, w int) (lsp, rsp string) {
	ns := StringWidth(s)
	if ns <= w {
		return lsp, rsp
	}
	nr := (w - ns) / 2
	nl := w - ns - nr
	lsp = Spaces(nl)
	rsp = Spaces(nr)
	return lsp, rsp
}
