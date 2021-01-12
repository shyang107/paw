package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func main() {
	a := []int{1, 2, 3}
	b := []float32{1., 2., 3.}
	c := []float64{1., 2., 3.}
	d := []string{"abc", "Abc", "aBc"}
	fmt.Printf("%#v, sum = %d\n", a, paw.Sum(a).(int))
	fmt.Printf("%#v, sum = %.2f\n", b, paw.Sum(b).(float32))
	fmt.Printf("%#v, sum = %.6f\n", c, paw.Sum(c).(float64))
	fmt.Printf("%#v, sum = %v\n", d, paw.Sum(d).(string))
	fmt.Println()
	fmt.Printf("%#v, min = %d\n", a, paw.Min(a).(int))
	fmt.Printf("%#v, min = %.2f\n", b, paw.Min(b).(float32))
	fmt.Printf("%#v, min = %.6f\n", c, paw.Min(c).(float64))
	fmt.Printf("%#v, min = %v\n", d, paw.Min(d).(string))
	fmt.Println()
	fmt.Printf("%#v, max = %d\n", a, paw.Max(a).(int))
	fmt.Printf("%#v, max = %.2f\n", b, paw.Max(b).(float32))
	fmt.Printf("%#v, max = %.6f\n", c, paw.Max(c).(float64))
	fmt.Printf("%#v,  max = %q\n", d, paw.Max(d).(string))
}
