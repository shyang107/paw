package _junk

import (
	"math/rand"
	"time"
)

// NewRand return a instance of
func NewRand() *rand.Rand {
	// return rand.New(rand.NewSource(time.Now().Unix()))
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Shuffle randomly shuffle the order of `slice`
func Shuffle(slice []interface{}) {
	r := NewRand()
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

// GenerateSlice generates a slice of size, size filled with random numbers
func GenerateSlice(size int) []int {
	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(999) - rand.Intn(999)
	}
	return slice
}

// CopySliceInt copy []int
func CopySliceInt(nums []int) []int {
	n := make([]int, len(nums), len(nums))
	copy(n, nums)
	return n
}

// ReverseSliceInt 反轉切片 nums 的 [i, j] 範圍
func ReverseSliceInt(nums []int, i, j int) {
	for i < j {
		nums[i], nums[j] = nums[j], nums[i]
		i++
		j--
	}
}

// // CombSort is a variant of the Bubble Sort in asscending order
// // 	BUG shrink 選擇不當有可能失敗, shrink >= 1.3
// func CombSort(nums []int, shrink float64) {
// 	var (
// 		n   = len(nums)
// 		gap = len(nums)
// 		// shrink  = 1.8
// 		swapped = true
// 	)

// 	for swapped {
// 		swapped = false
// 		gap = int(float64(gap) / shrink)
// 		if gap < 1 {
// 			gap = 1
// 		}
// 		for i := 0; i+gap < n; i++ {
// 			if nums[i] > nums[i+gap] {
// 				nums[i+gap], nums[i] = nums[i], nums[i+gap]
// 				swapped = true
// 			}
// 		}
// 	}
// }

// // CombSortFunc is a variant of the Bubble Sort in asscending order
// // 	f(a,b) : true for exchange, false not
// // Example:
// // 	升冪:  f func(a, b int) bool { rerturn a > b}
// // 	降冪:  f func(a, b int) bool { rerturn a < b}
// // 	BUG shrink 選擇不當有可能失敗, shrink >= 1.3
// func CombSortFunc(nums []int, shrink float64, f func(a, b int) bool) {
// 	var (
// 		n   = len(nums)
// 		gap = len(nums)
// 		// shrink  = 1.8
// 		swapped = true
// 	)

// 	for swapped {
// 		swapped = false
// 		gap = int(float64(gap) / shrink)
// 		if gap < 1 {
// 			gap = 1
// 		}
// 		for i := 0; i+gap < n; i++ {
// 			if f(nums[i], nums[i+gap]) {
// 				nums[i+gap], nums[i] = nums[i], nums[i+gap]
// 				swapped = true
// 			}
// 		}
// 	}
// }

// // MergeSort is a Divide and Conquer algorithm. Meaning, the algorithm splits an input into various pieces, sorts them and then merges them back together.
// func MergeSort(items []int) []int {
// 	var num = len(items)

// 	if num == 1 {
// 		return items
// 	}

// 	middle := int(num / 2)
// 	var (
// 		left  = make([]int, middle)
// 		right = make([]int, num-middle)
// 	)
// 	for i := 0; i < num; i++ {
// 		if i < middle {
// 			left[i] = items[i]
// 		} else {
// 			right[i-middle] = items[i]
// 		}
// 	}

// 	return merge(MergeSort(left), MergeSort(right))
// }

// func merge(left, right []int) (result []int) {
// 	result = make([]int, len(left)+len(right))

// 	i := 0
// 	for len(left) > 0 && len(right) > 0 {
// 		if left[0] < right[0] {
// 			result[i] = left[0]
// 			left = left[1:]
// 		} else {
// 			result[i] = right[0]
// 			right = right[1:]
// 		}
// 		i++
// 	}

// 	for j := 0; j < len(left); j++ {
// 		result[i] = left[j]
// 		i++
// 	}
// 	for j := 0; j < len(right); j++ {
// 		result[i] = right[j]
// 		i++
// 	}

// 	return
// }

// func bruteForceHelper(nums []int, n int, ans *[][]int) {
// 	if n == 1 {
// 		*ans = append(*ans, CopySliceInt(nums))
// 		return
// 	}

// 	for i := 0; i < n; i++ {
// 		nums[i], nums[n-1] = nums[n-1], nums[i]
// 		bruteForceHelper(nums, n-1, ans)
// 		nums[i], nums[n-1] = nums[n-1], nums[i]
// 	}
// }

// // BruteForce 通过暴力法生成一个序列的全部排列
// func BruteForce(nums []int) [][]int {
// 	ans := make([][]int, 0, len(nums))
// 	bruteForceHelper(nums, len(nums), &ans)
// 	return ans
// }
