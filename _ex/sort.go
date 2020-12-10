package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

// SelectionSort 選擇排序
func SelectionSort(n []int) {
	fmt.Println("SelectionSort\n", n)
	// count := 0
	// for i := 0; i < len(n); i++ {
	// 	minIndex := i
	// 	for j := i + 1; j < len(n); j++ {
	// 		if n[minIndex] > n[j] {
	// 			minIndex = j
	// 		}
	// 	}
	// 	n[i], n[minIndex] = n[minIndex], n[i]
	// 	count++
	// 	fmt.Println(count, n)
	// }
	paw.SelectionSort(n)
	fmt.Println(n)
	paw.SelectionSortFunc(n, func(a, b int) bool { return a < b })
	fmt.Println(n)

}

// InsertionSort 插入排序
func InsertionSort(n []int) {
	fmt.Println("InsertionSort\n", n)
	// count := 0
	// i := 1
	// for i < len(a) {
	// 	j := i
	// 	for j >= 1 && a[j] < a[j-1] {
	// 		a[j-1], a[j] = a[j], a[j-1]
	// 		count++
	// 		fmt.Println(count, a)
	// 		j--
	// 	}
	// 	i++
	// }
	paw.InsertionSort(n)
	fmt.Println(n)
	paw.InsertionSortFunc(n, func(a, b int) bool { return a < b })
	fmt.Println(n)
}
