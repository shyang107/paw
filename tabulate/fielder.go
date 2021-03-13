package tabulate

import "github.com/fatih/color"

type Fielder interface {
	Name() string
	Value() interface{}
	Width() int
	AlignS() string
	Color() *color.Color
}
