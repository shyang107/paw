package paw

import (
	"strings"

	"github.com/keakon/golog"
)

const (
// // InfoLevel is golog's level
// InfoLevel = golog.InfoLevel
// // WarnLeveL is golog's level
// WarnLeveL = golog.WarnLevel
// // DebugLevel is golog's level
// DebugLevel = golog.DebugLevel
// // CritLevel is golog's level
// CritLevel = golog.CritLevel
// // ErrorLevel is golog's level
// ErrorLevel = golog.ErrorLevel
)

var (
	// Log is log instance
	Log = golog.NewStderrLogger()
	sb  = strings.Builder{}
)

func init() {

}

// // SetLogLevel set the level of Log to `lv`
// //
// // 	`lv`: golog.LeveL (such as `InfoLevel`, `WarnLeveL`, `DebugLevel`, `CritLevel`, `ErrorLevel`)
// func SetLogLevel(lv golog.Level) {
// 	Log = golog.NewLogger(lv)
// }
