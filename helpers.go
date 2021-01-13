package paw

import (
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"

	"github.com/shyang107/paw/cast"
	// log "github.com/sirupsen/logrus"
)

// GetFuncName get the func name
func GetFuncName(level interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(level).Pointer()).Name()
}

// GetDate get the current date and return string. For example,
//	if time.Now().Date() is `2020 November 27`
//	then GetDate() return `Nov 27, 2020`
func GetDate() string {
	// y, m, d := time.Now().Date()
	// ms := string(([]byte(m.String()))[0:3])
	// return fmt.Sprintf("%s. %d, %d", ms, d, y)
	return time.Now().Format("Jan 2, 2006")
	//	GetDate() return carbon.Now().Format(carbon.FormattedDateFormat)
	//		carbon.FormattedDateFormat = "Jan 2, 2006"
	// return carbon.Now().Format(carbon.FormattedDateFormat)

}

// Max return the maximum value of `x`, the type of x is one of []int, []float32, []float64, []string.
// If x is empty, and type is one of []int, []float32 and []float64, will return 0.
// If x is empty, and type is one of []string, will return "".
func Max(x interface{}) interface{} {
	var r interface{}
	switch x.(type) {
	case []int:
		return MaxInts(x.([]int)...)
	case []float32:
		return MaxFloat32s(x.([]float32)...)
	case []float64:
		return MaxFloat64s(x.([]float64)...)
	case []string:
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
	switch x.(type) {
	case []int:
		return MinInts(x.([]int)...)
	case []float32:
		return MinFloat32s(x.([]float32)...)
	case []float64:
		return MinFloat64s(x.([]float64)...)
	case []string:
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
	switch x.(type) {
	case []int:
		return SumInts(x.([]int)...)
	case []float32:
		return SumFloat32s(x.([]float32)...)
	case []float64:
		return SumFloat64s(x.([]float64)...)
	case []string:
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
		// Error.Println("run stty, err: ", err)
		return 38, 100
	}
	size := Split(TrimSuffix(string(out), "\n"), " ")
	height = cast.ToInt(size[0])
	width = cast.ToInt(size[1])
	return height, width
}

// IndexOfString gets the index at which the first occurrence of a string value is found in array or return -1
// if the value cannot be found
func IndexOfString(a []string, x string) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

func indexOf(n int, f func(int) bool) int {
	for i := 0; i < n; i++ {
		if f(i) {
			return i
		}
	}
	return -1
}

// IndexOfInt gets the index at which the first occurrence of an int value is found in array or return -1
// if the value cannot be found
func IndexOfInt(a []int, x int) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfInt32 gets the index at which the first occurrence of an int32 value is found in array or return -1
// if the value cannot be found
func IndexOfInt32(a []int32, x int32) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfInt64 gets the index at which the first occurrence of an int64 value is found in array or return -1
// if the value cannot be found
func IndexOfInt64(a []int64, x int64) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfUInt gets the index at which the first occurrence of an uint value is found in array or return -1
// if the value cannot be found
func IndexOfUInt(a []uint, x uint) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfUInt32 gets the index at which the first occurrence of an uint32 value is found in array or return -1
// if the value cannot be found
func IndexOfUInt32(a []uint32, x uint32) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfUInt64 gets the index at which the first occurrence of an uint64 value is found in array or return -1
// if the value cannot be found
func IndexOfUInt64(a []uint64, x uint64) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// IndexOfFloat64 gets the index at which the first occurrence of an float64 value is found in array or return -1
// if the value cannot be found
func IndexOfFloat64(a []float64, x float64) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

func lastIndexOf(n int, f func(int) bool) int {
	for i := n - 1; i >= 0; i-- {
		if f(i) {
			return i
		}
	}
	return -1
}

// LastIndexOfInt gets the index at which the first occurrence of an int value is found in array or return -1
// if the value cannot be found
func LastIndexOfInt(a []int, x int) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfInt32 gets the index at which the first occurrence of an int32 value is found in array or return -1
// if the value cannot be found
func LastIndexOfInt32(a []int32, x int32) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfInt64 gets the index at which the first occurrence of an int64 value is found in array or return -1
// if the value cannot be found
func LastIndexOfInt64(a []int64, x int64) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfUInt gets the index at which the first occurrence of an uint value is found in array or return -1
// if the value cannot be found
func LastIndexOfUInt(a []uint, x uint) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfUInt32 gets the index at which the first occurrence of an uint32 value is found in array or return -1
// if the value cannot be found
func LastIndexOfUInt32(a []uint32, x uint32) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfUInt64 gets the index at which the first occurrence of an uint64 value is found in array or return -1
// if the value cannot be found
func LastIndexOfUInt64(a []uint64, x uint64) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfFloat64 gets the index at which the first occurrence of an float64 value is found in array or return -1
// if the value cannot be found
func LastIndexOfFloat64(a []float64, x float64) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfFloat32 gets the index at which the first occurrence of an float32 value is found in array or return -1
// if the value cannot be found
func LastIndexOfFloat32(a []float32, x float32) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOfString gets the index at which the first occurrence of a string value is found in array or return -1
// if the value cannot be found
func LastIndexOfString(a []string, x string) int {
	return lastIndexOf(len(a), func(i int) bool { return a[i] == x })
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

// ShuffleInt creates an array of int shuffled values using Fisher–Yates algorithm
func ShuffleInt(a []int) []int {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleInt32 creates an array of int32 shuffled values using Fisher–Yates algorithm
func ShuffleInt32(a []int32) []int32 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleInt64 creates an array of int64 shuffled values using Fisher–Yates algorithm
func ShuffleInt64(a []int64) []int64 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleUInt creates an array of int shuffled values using Fisher–Yates algorithm
func ShuffleUInt(a []uint) []uint {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleUInt32 creates an array of uint32 shuffled values using Fisher–Yates algorithm
func ShuffleUInt32(a []uint32) []uint32 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleUInt64 creates an array of uint64 shuffled values using Fisher–Yates algorithm
func ShuffleUInt64(a []uint64) []uint64 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleString creates an array of string shuffled values using Fisher–Yates algorithm
func ShuffleString(a []string) []string {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleFloat32 creates an array of float32 shuffled values using Fisher–Yates algorithm
func ShuffleFloat32(a []float32) []float32 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// ShuffleFloat64 creates an array of float64 shuffled values using Fisher–Yates algorithm
func ShuffleFloat64(a []float64) []float64 {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

// RandomInt generates a random int, based on a min and max values
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// ContainsInt returns true if an int is present in a iteratee.
func ContainsInt(s []int, v int) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsInt32 returns true if an int32 is present in a iteratee.
func ContainsInt32(s []int32, v int32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsInt64 returns true if an int64 is present in a iteratee.
func ContainsInt64(s []int64, v int64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt returns true if an uint is present in a iteratee.
func ContainsUInt(s []uint, v uint) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt32 returns true if an uint32 is present in a iteratee.
func ContainsUInt32(s []uint32, v uint32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt64 returns true if an uint64 is present in a iteratee.
func ContainsUInt64(s []uint64, v uint64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsString returns true if a string is present in a iteratee.
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsFloat32 returns true if a float32 is present in a iteratee.
func ContainsFloat32(s []float32, v float32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsFloat64 returns true if a float64 is present in a iteratee.
func ContainsFloat64(s []float64, v float64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ReverseStrings reverses an array of string
func ReverseStrings(s []string) []string {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseInt reverses an array of int
func ReverseInt(s []int) []int {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseInt32 reverses an array of int32
func ReverseInt32(s []int32) []int32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseInt64 reverses an array of int64
func ReverseInt64(s []int64) []int64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt reverses an array of int
func ReverseUInt(s []uint) []uint {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt32 reverses an array of uint32
func ReverseUInt32(s []uint32) []uint32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt64 reverses an array of uint64
func ReverseUInt64(s []uint64) []uint64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseFloat64 reverses an array of float64
func ReverseFloat64(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseFloat32 reverses an array of float32
func ReverseFloat32(s []float32) []float32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseString reverses a string
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
