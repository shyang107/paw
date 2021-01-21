package main

import (
	"fmt"
	"reflect"
)

type tstruct struct {
	a int
}
type tslice []int
type tint int

func main() {
	var (
		a         = make([]int, 1)
		b         = make([]string, 2)
		c         = make([]tstruct, 3)
		d  tslice = make([]int, 4)
		e         = &d
		va        = reflect.ValueOf(a)
		vb        = reflect.ValueOf(b)
		vc        = reflect.ValueOf(c)
		vd        = reflect.ValueOf(d)
		ve        = reflect.ValueOf(e)
	)
	for idx := -1; idx < 6; idx++ {
		if err := CheckIndex(a, idx); err != nil {
			fmt.Println(err, va.Type(), va.Kind())
		} else {
			fmt.Println("idx =", idx, "is in range of slice", va.Type(), va.Kind())
		}
		if err := CheckIndex(b, idx); err != nil {
			fmt.Println(err, vb.Type(), vb.Kind())
		} else {
			fmt.Println("idx =", idx, "is in range of slice", vb.Type(), vb.Kind())
		}
		if err := CheckIndex(c, idx); err != nil {
			fmt.Println(err, vc.Type(), vc.Kind())
		} else {
			fmt.Println("idx =", idx, "is in range of slice", vc.Type(), vc.Kind())
		}
		if err := CheckIndex(d, idx); err != nil {
			fmt.Println(err, vd.Type(), vd.Kind())
		} else {
			fmt.Println("idx =", idx, "is in range of slice", vd.Type(), vd.Kind())
		}
		if err := CheckIndex(e, idx); err != nil {
			fmt.Println(err, ve.Type(), ve.Kind())
		} else {
			fmt.Println("idx =", idx, "is in range of slice", ve.Type(), ve.Kind())
		}
	}
	idx := 2
	v := reflect.ValueOf(tint(1))
	if err := CheckIndex(1, idx); err != nil {
		fmt.Println(err, v.Type(), v.Kind())
	} else {
		fmt.Println("idx =", idx, "is in range of slice", v.Type(), v.Kind())
	}
	v = reflect.ValueOf("1")
	if err := CheckIndex("1", idx); err != nil {
		fmt.Println(err, v.Type(), v.Kind())
	} else {
		fmt.Println("idx =", idx, "is in range of slice", v.Type(), v.Kind())
	}
}

func CheckIndex(slice interface{}, idx int) error {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("CheckIndex: expected slice type, found %q", v.Kind().String())
	}
	count := v.Len()
	if idx < 0 || idx > count-1 {
		return fmt.Errorf("CheckIndex: slice range [%d, %d), idx is %d", 0, count, idx)
	}
	return nil
}
