package paw

import (
	"math/rand"
	"time"
)

// NewRand return a instance of
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().Unix()))
}

// Shuffle randomly shuffle the order of `slice`
func Shuffle(slice []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

// SelectionSort 選擇排序 (升冪) (較慢)
func SelectionSort(n []int) {
	for i := 0; i < len(n); i++ {
		minIndex := i
		for j := i + 1; j < len(n); j++ {
			if n[minIndex] > n[j] {
				minIndex = j
			}
		}
		n[i], n[minIndex] = n[minIndex], n[i]
	}
}

// SelectionSortFunc 選擇排序 (升冪) (較慢)
// 	f(a,b) : true for exchange, false not
// Example:
// 	升冪:  f func(a, b int) bool { rerturn a > b}
// 	降冪:  f func(a, b int) bool { rerturn a < b}
func SelectionSortFunc(n []int, f func(a, b int) bool) {
	for i := 0; i < len(n); i++ {
		minIndex := i
		for j := i + 1; j < len(n); j++ {
			if f(n[minIndex], n[j]) {
				minIndex = j
			}
		}
		n[i], n[minIndex] = n[minIndex], n[i]
	}
}

// InsertionSort 插入排序 (升冪) (較快)
func InsertionSort(n []int) {
	i := 1
	for i < len(n) {
		j := i
		for j >= 1 && n[j-1] > n[j] {
			n[j-1], n[j] = n[j], n[j-1]
			j--
		}
		i++
	}
}

// InsertionSortFunc 插入排序 (升冪) (較快)
// 	f(a,b) : true for exchange, false not
// Example:
// 	升冪:  f func(a, b int) bool { rerturn a > b}
// 	降冪:  f func(a, b int) bool { rerturn a < b}
func InsertionSortFunc(n []int, f func(a, b int) bool) {
	i := 1
	for i < len(n) {
		j := i
		for j >= 1 && f(n[j-1], n[j]) {
			n[j-1], n[j] = n[j], n[j-1]
			j--
		}
		i++
	}
}
