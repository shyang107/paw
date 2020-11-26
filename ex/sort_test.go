package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shyang107/paw"
)

const (
	N = 50000
)

// var r = rand.New(rand.NewSource(time.Now().Unix()))
var r = paw.NewRand()

func GetRandomNums(n int) []int {
	// nums := []int{}
	// for i := 0; i < N; i++ {
	// 	nums = append(nums, r.Intn(N))
	// }
	// return nums
	return paw.GenerateSlice(n)
}

// func BenchmarkParallelSort(b *testing.B) {
// 	nums := GetRandomNums(N)
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			sort.Ints(nums)
// 		}
// 	})
// }
func BenchmarkSort(b *testing.B) {
	nums := GetRandomNums(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort.Ints(nums)
	}
}

// func BenchmarkParallelSelctionSort(b *testing.B) {
// 	nums := GetRandomNums(N)
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			paw.SelectionSort(nums)
// 		}
// 	})
// }
func BenchmarkSelectionSort(b *testing.B) {
	nums := GetRandomNums(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paw.SelectionSort(nums)
	}
}

// func BenchmarkParallelInsertionSort(b *testing.B) {
// 	nums := GetRandomNums(N)
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			paw.InsertionSort(nums)
// 		}
// 	})
// }
func BenchmarkInsertionSort(b *testing.B) {
	nums := GetRandomNums(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paw.InsertionSort(nums)
	}
}

// func BenchmarkCombSort(b *testing.B) {
// 	nums := GetRandomNums(N)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		paw.CombSort(nums)
// 	}
// }

func TestInsertionSort(t *testing.T) {
	var actual = []int{1, 39, 2, 9, 7, 54, 11}
	var expected = []int{1, 2, 7, 9, 11, 39, 54}
	paw.InsertionSort(actual)
	assert.Equal(t, expected, actual)
}
func TestSelectionSort(t *testing.T) {
	var actual = []int{1, 39, 2, 9, 7, 54, 11}
	var expected = []int{1, 2, 7, 9, 11, 39, 54}
	paw.SelectionSort(actual)
	assert.Equal(t, expected, actual)
}

var f = func(a, b int) bool { return a < b }

func TestInsertionSortFunc(t *testing.T) {
	var actual = []int{1, 39, 2, 9, 7, 54, 11}
	var expected = []int{54, 39, 11, 9, 7, 2, 1}
	paw.InsertionSortFunc(actual, f)
	assert.Equal(t, expected, actual)
}
func TestSelectionSortFunc(t *testing.T) {
	var actual = []int{1, 39, 2, 9, 7, 54, 11}
	var expected = []int{54, 39, 11, 9, 7, 2, 1}
	paw.SelectionSortFunc(actual, f)
	assert.Equal(t, expected, actual)
}
