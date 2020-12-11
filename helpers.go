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
