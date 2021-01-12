package paw

import (
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"

	"github.com/spf13/cast"
	"github.com/uniplaces/carbon"
	// log "github.com/sirupsen/logrus"
)

// GetFuncName get the func name
func GetFuncName(level interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(level).Pointer()).Name()
}

// GetDate get the current date and return string. For example,
//	GetDate() return carbon.Now().Format(carbon.FormattedDateFormat)
//	if time.Now().Date() is `2020 November 27`
//		carbon.FormattedDateFormat = "Jan 2, 2006"
//	then GetDate() return `Nov 27, 2020`
func GetDate() string {
	// y, m, d := time.Now().Date()
	// ms := string(([]byte(m.String()))[0:3])
	// return fmt.Sprintf("%s. %d, %d", ms, d, y)
	// return time.Now().Format("Jan. 2, 2006")
	return carbon.Now().Format(carbon.FormattedDateFormat)

}

// Max return the maximum value of `x`, the type of x is one of []int, []float32, []float64, []string.
// If x is empty, and type is one of []int, []float32 and []float64, will return 0.
// If x is empty, and type is one of []string, will return "".
func Max(x interface{}) interface{} {
	var r interface{}
	switch reflect.TypeOf(x) {
	case reflect.TypeOf([]int{}):
		return MaxInts(x.([]int)...)
	case reflect.TypeOf([]float32{}):
		return MaxFloat32s(x.([]float32)...)
	case reflect.TypeOf([]float64{}):
		return MaxFloat64s(x.([]float64)...)
	case reflect.TypeOf([]string{}):
		return MaxStrings(x.([]string)...)
	default:
		return nil
	}
	return r
}

// MaxInt will return maximum value of `i` and `j`
func MaxInt(i, j int) int {
	return MaxInts(i, j)
}

// MaxInts will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxInts(x ...int) int {
	if len(x) == 0 {
		return 0
	}
	var m int = x[0]
	for _, j := range x {
		if m < j {
			m = j
		}
	}

	return m
}

// MaxFloat32s will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxFloat32s(x ...float32) float32 {
	if len(x) == 0 {
		return 0
	}
	var m float32 = x[0]
	for _, j := range x {
		if m < j {
			m = j
		}
	}

	return m
}

// MaxFloat64s will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxFloat64s(x ...float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var m float64 = x[0]
	for _, j := range x {
		if m < j {
			m = j
		}
	}

	return m
}

// MaxStrings will return maximum value of `x...`.
// Compares iterms using code point.
func MaxStrings(x ...string) string {
	if len(x) == 0 {
		return ""
	}
	var m string = x[0]
	for _, v := range x {
		if m < v {
			m = v
		}
	}

	return m
}

// Min return the minimum value of `x`, the type of x is one of []int, []float32, []float64, []string.
// If x is empty, and type is one of []int, []float32 and []float64, will return 0.
// If x is empty, and type is one of []string, will return "".
func Min(x interface{}) interface{} {
	var r interface{}
	switch reflect.TypeOf(x) {
	case reflect.TypeOf([]int{}):
		return MinInts(x.([]int)...)
	case reflect.TypeOf([]float32{}):
		return MinFloat32s(x.([]float32)...)
	case reflect.TypeOf([]float64{}):
		return MinFloat64s(x.([]float64)...)
	case reflect.TypeOf([]string{}):
		return MinStrings(x.([]string)...)
	default:
		return nil
	}
	return r
}

// MinInt will return minimum value of `i` and `j`
func MinInt(i, j int) int {
	return MinInts(i, j)
}

// MinInts will return minimum value of `x...`.
// If x is empty, and return 0.
func MinInts(x ...int) int {
	if len(x) == 0 {
		return 0
	}
	var m int = x[0]
	for _, j := range x {
		if m > j {
			m = j
		}
	}

	return m
}

// MinFloat32s will return minimum value of `x...`.
// If x is empty, and return 0.
func MinFloat32s(x ...float32) float32 {
	if len(x) == 0 {
		return 0
	}
	var m float32 = x[0]
	for _, j := range x {
		if m > j {
			m = j
		}
	}

	return m
}

// MinFloat64s will return minimum value of `x...`.
// If x is empty, and return 0.
func MinFloat64s(x ...float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var m float64 = x[0]
	for _, j := range x {
		if m > j {
			m = j
		}
	}

	return m
}

// MinStrings will return minimum value of `x...`.
// Compares iterms using code point
// If x is empty, and return "".
func MinStrings(x ...string) string {
	if len(x) == 0 {
		return ""
	}
	var m string = x[0]
	for _, v := range x {
		if m > v {
			m = v
		}
	}

	return m
}

// Sum will return summation of []int, []float32, []flot64, []string.
// If the type of `x` is not one of []int, []float32, []flot64 and []string, then return `nil`
// If `x` is []string, will return concatenation using `strings.Join(a,"")`.
func Sum(x interface{}) interface{} {
	var r interface{}
	switch reflect.TypeOf(x) {
	case reflect.TypeOf([]int{}):
		return SumInts(x.([]int)...)
	case reflect.TypeOf([]float32{}):
		return SumFloat32s(x.([]float32)...)
	case reflect.TypeOf([]float64{}):
		return SumFloat64s(x.([]float64)...)
	case reflect.TypeOf([]string{}):
		return SumStrings(x.([]string)...)
	default:
		return nil
	}
	return r
}

// SumInts will return summation of []int
func SumInts(a ...int) int {
	var s int = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumFloat32s will return summation of []float32
func SumFloat32s(a ...float32) float32 {
	var s float32 = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumFloat64s will return summation of []float64
func SumFloat64s(a ...float64) float64 {
	var s float64 = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumStrings will return concatenation of []string using strings.Join(a, "")
func SumStrings(a ...string) string {
	return Join(a, "")
}

// NewRand return a instance of
func NewRand() *rand.Rand {
	// return rand.New(rand.NewSource(time.Now().Unix()))
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Shuffle randomly shuffle the order of `slice`
func Shuffle(slice []interface{}) {
	r := NewRand()
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

// PaddingString add pad-prefix in every line of string
func PaddingString(s string, pad string) string {
	if !Contains(s, "\n") {
		return pad + s
	}
	ss := Split(s, "\n")
	// sb := strings.Builder{}
	sb := ""
	for i := 0; i < len(ss)-1; i++ {
		sb += pad + ss[i] + "\n"
	}
	sb += pad + ss[len(ss)-1]
	return sb
}

// PaddingBytes add pad-prefix in every line('\n') of []byte
func PaddingBytes(bytes []byte, pad string) []byte {
	b := make([]byte, len(bytes))
	b = append(b, pad...)
	for _, v := range bytes {
		b = append(b, v)
		if v == '\n' {
			b = append(b, pad...)
		}
	}
	return b
}

// GetTerminalSize get size of console using `stty size`
func GetTerminalSize() (height, width int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// log.Fatal(err)
		Error.Println("run stty, err: ", err)
		return 38, 100
	}
	size := Split(TrimSuffix(string(out), "\n"), " ")
	height = cast.ToInt(size[0])
	width = cast.ToInt(size[1])
	return height, width
}
