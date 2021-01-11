package paw

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"

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

// MaxInt will return maximum value of `i` and `j`
func MaxInt(i, j int) int {
	return MaxInts(i, j)
}

// MaxInts will return maximum value of `i` and `js...`
func MaxInts(i int, js ...int) int {
	if len(js) == 0 {
		return i
	}
	m := i
	for _, j := range js {
		if m < j {
			m = j
		}
	}

	return m
}

// MinInt will return minimum value of `i` and `j`
func MinInt(i, j int) int {
	return MinInts(i, j)
}

// MinInts will return minimum value of `i` and `js...`
func MinInts(i int, js ...int) int {
	if len(js) == 0 {
		return i
	}
	m := i
	for _, j := range js {
		if m > j {
			m = j
		}
	}

	return m
}

// SumInts will return summation of []int
func SumInts(a ...int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
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
