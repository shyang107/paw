package paw

import (
	"reflect"
	"runtime"

	"github.com/uniplaces/carbon"
	// log "github.com/sirupsen/logrus"
)

// GetFuncName get the func name
func GetFuncName(level interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(level).Pointer()).Name()
}

// GetDate get the current date and return string
func GetDate() string {
	// y, m, d := time.Now().Date()
	// ms := string(([]byte(m.String()))[0:3])
	// return fmt.Sprintf("%s. %d, %d", ms, d, y)
	// return time.Now().Format("Jan. 2, 2006")
	return carbon.Now().Format(carbon.FormattedDateFormat)

}
