package paw

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// LineCount counts the number of '\n' for reader `r`
func LineCount(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

// FileLineCount counts the number of '\n' for file f
// 	`f` could be `gzip` file or plain text file
//
// modify from "github.com/liuzl/goutil"
func FileLineCount(f string) (int, error) {
	if strings.HasSuffix(strings.ToLower(f), ".gz") {
		fr, err := os.Open(f)
		if err != nil {
			return 0, err
		}
		defer fr.Close()
		r, err := gzip.NewReader(fr)
		if err != nil {
			return 0, err
		}
		return LineCount(r)
	}
	r, err := os.Open(f)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	return LineCount(r)
}

// ForEachLine higher order function that processes each line of text by callback function.
// The last non-empty line of input will be processed even if it has no newline.
// 	`br` : read from `br` reader
// 	`callback` : the function used to treatment the each line from `br`
//
// modify from "github.com/liuzl/goutil"
func ForEachLine(br *bufio.Reader, callback func(string) error) error {
	stop := false
	for {
		if stop {
			break
		}
		line, err := br.ReadString('\n')
		if err == io.EOF {
			stop = true
		} else if err != nil {
			return err
		}
		line = strings.TrimSuffix(line, "\n")
		if line == "" {
			if !stop {
				if err = callback(line); err != nil {
					return err
				}
			}
			continue
		}
		if err = callback(line); err != nil {
			return err
		}
	}
	return nil
}

// AppendToFile append string `s` to file `fileName`
func AppendToFile(fileName string, s string) error {
	var file *os.File
	var err error
	if IsFileExist(fileName) {
		file, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(fileName)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(s); err != nil {
		return err
	}
	return nil
}

// Major returns the major component of a Darwin device number.
func Major(dev uint64) uint32 {
	// copy from https://golang.org/src/archive/tar/stat_unix.go
	return uint32((dev >> 24) & 0xff)
}

// Minor returns the minor component of a Darwin device number.
func Minor(dev uint64) uint32 {
	// copy from https://golang.org/src/archive/tar/stat_unix.go
	return uint32(dev & 0xffffff)
}

// Mkdev returns a Darwin device number generated from the given major and minor
// components.
func Mkdev(major, minor uint32) uint64 {
	return (uint64(major) << 24) | uint64(minor)
}

// DevNumber returns the major and minor component of a Darwin device number.
func DevNumber(dev uint64) (uint32, uint32) {
	return Major(dev), Minor(dev)
}

// GenMd5 will generate md5 string of file (path)
func GenMd5(fullPath string) string {
	f, err := os.Open(fullPath)
	if err != nil {
		return err.Error()
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
	// return hex.EncodeToString(h.Sum(nil))
}

// GenMd5sh will generate md5 string of file (path)
func GenMd5sh(fullPath string) string {
	fi, err := os.Stat(fullPath)
	if err != nil {
		return err.Error()
	}
	out, err := execOutput(fmt.Sprintf("md5 %q", fullPath))
	if err != nil || fi.IsDir() {
		return "-"
	}
	md5s := strings.Split(strings.TrimSpace(out), " = ")
	return md5s[1]
}

var execOutput = func(cmd string) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()

	return string(out), err
}

// // GenMd5sh will generate md5 string of file (path)
// func GenMd5sh(fullPath string) string {
// 	fi, err := os.Stat(fullPath)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	out, err := execOutput(fmt.Sprintf("md5 %q", fullPath))
// 	if err != nil || fi.IsDir() {
// 		return "-"
// 	}
// 	md5s := strings.Split(strings.TrimSpace(out), " = ")
// 	return md5s[1]
// }

// var execOutput = func(cmd string) (string, error) {
// 	out, err := exec.Command("/bin/sh", "-c", cmd).Output()

// 	return string(out), err
// }
