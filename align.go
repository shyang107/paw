package paw

// Align is id that indicate alignment of head-column
type Align int

const (
	// AlignLeft align left
	AlignLeft Align = iota + 1
	// AlignCenter align center
	AlignCenter
	// AlignRight align right
	AlignRight
)

func (a Align) String() string {
	switch a {
	case AlignLeft:
		return "left"
	case AlignCenter:
		return "center"
	case AlignRight:
		return "right"
	default:
		return "Unknown"
	}
}
