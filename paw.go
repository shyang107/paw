package paw

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// TODO version message

// app is the main structure of a cli application. It is recommended that
// an app be created with the cli.NewApp() function
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Version of the program
	Version string
	// Description of the program
	Description string
	// Compilation date
	Compiled time.Time
	// List of all authors who contributed
	Authors []*Author
	// Copyright of the binary if any
	Copyright string
	// Reader reader to write input to (useful for tests)
	Reader io.Reader
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
	// Tags recording main changes of version
	Tags map[string][]string
}

// Author represents someone who has contributed to a cli project.
type Author struct {
	Name  string // The Authors name
	Email string // The Authors email
}

// String makes Author comply to the Stringer interface, to allow an easy print in the templating process
func (a *Author) String() string {
	e := ""
	if a.Email != "" {
		e = " <" + a.Email + ">"
	}

	return fmt.Sprintf("%v%v", a.Name, e)
}

// Tries to find out when this binary was compiled.
// Returns the current time if it fails to find it.
func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

// NewApp creates a new cli Application with some reasonable defaults for Name,
// Usage, Version and Action.
func newApp() *App {
	return &App{
		Name:        filepath.Base(os.Args[0]),
		Version:     "",
		Description: "A new paw utility package",
		Compiled:    compileTime(),
		// Authors []*Author,
		// Copyright string
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Tags:      tags,
	}
}

var (
	app = &App{
		Name:        filepath.Base(os.Args[0]),
		Version:     "v0.0.7.5",
		Description: "A new paw utility package",
		Compiled:    compileTime(), //cast.ToTime("2021-02-22")
		Authors: []*Author{
			{
				Name:  "shyang",
				Email: "shyang107@gmail.com",
			},
		},
		// Copyright string
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
	}

	// Tags stores abstracts of every version to solving main problems
	tags = map[string][]string{
		"v0.0.7.5": { // "2021.2.22",
			"FindFiles using ReadDir",
			"fixed git status of dir (reflect git status of sub-files) and sub-files",
			"add md5 fields",
			"use goroutine to get md5",
		},
	}
)
