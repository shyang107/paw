package paw

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var w io.Writer = os.Stdout

// Head is the struct that is head part of report
type Head struct {
	Columns       []string
	Widths        []int
	Aligns        []Align
	TopSepChar    string
	BottomSepChar string
}

var hd *Head

// NewHead is the constructor of type Head
func NewHead() *Head {
	hd = &Head{
		Columns:       []string{},
		Widths:        []int{},
		Aligns:        []Align{},
		TopSepChar:    "=",
		BottomSepChar: "-"}
	return hd
}

// Align is id that indicate alignment of head-column
type Align int

const (
	// LEFT is that align left
	LEFT Align = iota
	// CENTER is that align center
	CENTER
	// RIGHT is that align right
	RIGHT
)

var sp = " "

// Print print Head string on stdout
func (h *Head) Print() error {
	if len(h.Columns) != len(h.Widths) {
		return errors.New("parameters: assignment mismatch")
	}
	hlen := h.GetHeadLength()
	// log.Println(hlen)
	top := strings.Repeat(h.TopSepChar, hlen)
	bottom := getBottom(h)
	row := getHeadRow(h)
	fmt.Fprintln(w, top)
	fmt.Fprintln(w, row)
	fmt.Fprintln(w, bottom)
	return nil
}

func getHeadRow(h *Head) string {
	r := ""
	n := len(h.Widths) - 1
	if len(h.Aligns) < n+1 {
		count := n + 1 - len(h.Aligns)
		for i := 0; i < count; i++ {
			h.Aligns = append(h.Aligns, LEFT)
		}
	}
	for i, lw := range h.Widths {
		s := h.Columns[i]
		lh := len(s)
		if lw >= lh {
			al := h.Aligns[i]
			switch al {
			case LEFT:
				r += s + strings.Repeat(sp, lw-lh)
			case RIGHT:
				r += strings.Repeat(sp, lw-lh) + s
			case CENTER:
				nr := (lw - lh) / 2
				nl := lw - lh - nr
				r += strings.Repeat(sp, nl) + s + strings.Repeat(sp, nr)
			}
		} else {
			r += string(([]byte(s))[0:lw])
		}
		if i < n {
			r += sp
		}
	}
	return r
}

func getBottom(h *Head) string {
	b := ""
	n := len(h.Widths) - 1
	for i, count := range h.Widths {
		if i < n {
			b += strings.Repeat(h.BottomSepChar, count) + sp
		} else {
			b += strings.Repeat(h.BottomSepChar, count)
		}
	}
	return b
}

// GetHeadLength calculate and return total length of head
func (h *Head) GetHeadLength() int {
	l := 0
	for _, lw := range h.Widths {
		l += lw
	}
	l += len(h.Widths) - 1
	return l
}

// PrintRow print out row values
func (h *Head) PrintRow(values ...string) error {
	if len(h.Widths) != len(values) {
		return errors.New("len: assignment mismatch")
	}
	n := len(h.Widths) - 1
	r := ""
	for i, lw := range h.Widths {
		s := values[i]
		lh := len(s)
		if lw >= lh {
			al := h.Aligns[i]
			switch al {
			case LEFT:
				r += s + strings.Repeat(sp, lw-lh)
			case RIGHT:
				r += strings.Repeat(sp, lw-lh) + s
			case CENTER:
				nr := (lw - lh) / 2
				nl := lw - lh - nr
				r += strings.Repeat(sp, nl) + s + strings.Repeat(sp, nr)
			}
		} else {
			r += string(([]byte(s))[0:lw])
		}
		if i < n {
			r += sp
		}
	}
	fmt.Fprintln(w, r)
	return nil
}
