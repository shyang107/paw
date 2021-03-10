package paw

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cast"
)

func Caller(skip int) string {
	if skip < 0 {
		skip = 0
	}
	pc, path, line, ok := runtime.Caller(skip + 1)
	if ok {
		function := runtime.FuncForPC(pc).Name()
		s := strings.Split(function, ".")
		funcName := s[len(s)-1]
		c := FileLSColor(path)
		base := filepath.Base(path)
		return Cdashp.Sprint(" from [") +
			c.Sprint(base) + Cdashp.Sprint(":") +
			Csnp.Sprint(line) + Cdashp.Sprint("][") +
			color.New(color.FgYellow).Sprint(funcName) + Cdashp.Sprint("]")
	}
	return fmt.Errorf("Caller(%d) failed, %s", skip, path).Error()
}

// GetFuncName get the func name
func GetFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
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

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// Max return the maximum value of `x`, the type of x is one of []int, []float32, []float64, []string.
// If x is empty, and type is one of []int, []float32 and []float64, will return 0.
// If x is empty, and type is one of []string, will return "".
func Max(x interface{}) interface{} {
	x = indirect(x)

	// var r interface{}
	switch x.(type) {
	case []int:
		return MaxIntA(x.([]int)...)
	case []float32:
		return MaxFloat32A(x.([]float32)...)
	case []float64:
		return MaxFloat64A(x.([]float64)...)
	case []string:
		return MaxStringA(x.([]string)...)
	default:
		return nil
	}
	// return r
}

// MaxInt will return maximum value of `i` and `j`
func MaxInt(i, j int) int {
	return MaxIntA(i, j)
}

// MaxIntA will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxIntA(x ...int) int {
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

// MaxFloat32A will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxFloat32A(x ...float32) float32 {
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

// MaxFloat64A will return maximum value of `x...`.
// If x is empty, and return 0.
func MaxFloat64A(x ...float64) float64 {
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

// MaxStringA will return maximum value of `x...`.
// Compares iterms using code point.
func MaxStringA(x ...string) string {
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
	x = indirect(x)

	// var r interface{}
	switch x.(type) {
	case []int:
		return MinIntA(x.([]int)...)
	case []float32:
		return MinFloat32A(x.([]float32)...)
	case []float64:
		return MinFloat64A(x.([]float64)...)
	case []string:
		return MinStringA(x.([]string)...)
	default:
		return nil
	}
	// return r
}

// MinInt will return minimum value of `i` and `j`
func MinInt(i, j int) int {
	return MinIntA(i, j)
}

// MinIntA will return minimum value of `x...`.
// If x is empty, and return 0.
func MinIntA(x ...int) int {
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

// MinFloat32A will return minimum value of `x...`.
// If x is empty, and return 0.
func MinFloat32A(x ...float32) float32 {
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

// MinFloat64A will return minimum value of `x...`.
// If x is empty, and return 0.
func MinFloat64A(x ...float64) float64 {
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

// MinStringA will return minimum value of `x...`.
// Compares iterms using code point
// If x is empty, and return "".
func MinStringA(x ...string) string {
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

// SumMap manipulates an iterate and sum it, like as SumMapE and if error occures, will return -1, but ignore error.
func SumMap(a interface{}, mapFunc func(idx int) int) int {
	s, _ := SumMapE(a, mapFunc)
	return s
}

// SumMapE manipulates an iterate and sum it, and if error occures, will return -1, error.
func SumMapE(a interface{}, mapFunc func(idx int) int) (int, error) {
	a = indirect(a)

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Slice {
		return -1, fmt.Errorf("CheckIndex: expected slice type, found %q", v.Kind().String())
	}

	count := v.Len()
	wd := 0
	for i := 0; i < count; i++ {
		wd += mapFunc(i)
	}

	return wd, nil
}

// Sum will return summation of []int, []float32, []flot64, []string.
// If the type of `x` is not one of []int, []float32, []flot64 and []string, then return `nil`
// If `x` is []string, will return concatenation using `strings.Join(a,"")`.
func Sum(x interface{}) interface{} {
	x = indirect(x)

	// var r interface{}
	switch x.(type) {
	case []int:
		return SumIntA(x.([]int)...)
	case []float32:
		return SumFloat32A(x.([]float32)...)
	case []float64:
		return SumFloat64A(x.([]float64)...)
	case []string:
		return SumStringA(x.([]string)...)
	default:
		return nil
	}
	// return r
}

// SumIntA will return summation of []int
func SumIntA(a ...int) int {
	var s int = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumFloat32A will return summation of []float32
func SumFloat32A(a ...float32) float32 {
	var s float32 = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumFloat64A will return summation of []float64
func SumFloat64A(a ...float64) float64 {
	var s float64 = 0
	for _, v := range a {
		s += v
	}
	return s
}

// SumStringA will return concatenation of []string using strings.Join(a, "")
func SumStringA(a ...string) string {
	return strings.Join(a, "")
}

// NewRand return a instance of
func NewRand() *rand.Rand {
	// return rand.New(rand.NewSource(time.Now().Unix()))
	return rand.New(rand.NewSource(time.Now().UnixNano()))
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
	size := strings.Split(strings.TrimSpace(string(out)), " ")
	height = cast.ToInt(size[0])
	width = cast.ToInt(size[1])
	return height, width
}

// IndexOf gets the index at which the first occurrence of an value is found in array or return -1.
// if the value cannot be found
func IndexOf(n int, f func(int) bool) int {
	return indexOf(n, f)
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

// IndexOfString gets the index at which the first occurrence of a string value is found in array or return -1
// if the value cannot be found
func IndexOfString(a []string, x string) int {
	return indexOf(len(a), func(i int) bool { return a[i] == x })
}

// LastIndexOf gets the index at which the first occurrence of an value is found in array or return -1.
// if the value cannot be found
func LastIndexOf(n int, f func(int) bool) int {
	return lastIndexOf(n, f)
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
		if strings.EqualFold(vv, v) {
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

// ReverseStringA reverses an array of string
func ReverseStringA(s []string) []string {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseIntA reverses an array of int
func ReverseIntA(s []int) []int {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseInt32A reverses an array of int32
func ReverseInt32A(s []int32) []int32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseInt64A reverses an array of int64
func ReverseInt64A(s []int64) []int64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt reverses an array of int
func ReverseUIntA(s []uint) []uint {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt32 reverses an array of uint32
func ReverseUInt32A(s []uint32) []uint32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseUInt64 reverses an array of uint64
func ReverseUInt64A(s []uint64) []uint64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseFloat64 reverses an array of float64
func ReverseFloat64A(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// ReverseFloat32 reverses an array of float32
func ReverseFloat32A(s []float32) []float32 {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
