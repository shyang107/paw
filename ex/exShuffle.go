package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func exShuffle() {
	s := []rune("abcdefg")
	slice := make([]interface{}, len(s))
	for i, val := range s {
		slice[i] = string(val)
	}
	fmt.Println(slice)
	for i := 0; i < 10; i++ {
		paw.Shuffle(slice)
		fmt.Println(slice)
	}
}
