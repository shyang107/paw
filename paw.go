package paw

import (
	"strings"

	"github.com/keakon/golog"
)

var (
	// Log is log instance
	Log = golog.NewStderrLogger()
	sb  = strings.Builder{}
)
