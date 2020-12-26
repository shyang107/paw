package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func checkType() {
	var winterfaces []interface{}
	stdout := os.Stdout
	stderr := os.Stderr
	buffer := new(bytes.Buffer)
	multiwriter := io.MultiWriter(stdout, buffer)
	winterfaces = append(winterfaces, stdout, stderr, buffer, multiwriter)
	buffer.WriteString("buffer")
	for _, v := range winterfaces {
		// t := reflect.TypeOf(v)
		// fmt.Println(t)
		switch vv := v.(type) {
		case *os.File:
			fmt.Printf("%v is type  *os.File %v.\n", v, vv)
		case *bytes.Buffer:
			fmt.Printf("%v is type  *bytes.Buffer %v.\n", v, vv)
		default:
			fmt.Printf("%v is type  {{writer}} %v.\n", v, vv)

		}
		if sv, ok := v.(*bytes.Buffer); ok {
			fmt.Printf("v implements *bytes.Buffer(): %s\n", sv.Bytes()) // note: sv, not v
		}
	}

}
