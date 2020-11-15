package paw

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRand(t *testing.T) {
	tests := []struct {
		name string
		want *rand.Rand
	}{
		{
			name: "NewRand",
			want: NewRand(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	s := []rune("abcdefg")
	slice := make([]interface{}, len(s))
	for i, val := range s {
		slice[i] = string(val)
	}
	type args struct {
		slice []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Shuffle",
			args: args{slice: slice},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Shuffle(tt.args.slice)
		})
	}
}

var (
	actual   = []int{1, 39, 2, 9, 7, 54, 11}
	expected = []int{1, 2, 7, 9, 11, 39, 54}
	sortfunc = func(a, b int) bool { return a < b }
)

type sortArgs struct {
	n []int
	w []int
}
type sortFuncArgs struct {
	sortArgs
	f func(a, b int) bool
}

type testSortCase struct {
	name string
	args sortArgs
}

type testSortFuncCase struct {
	name string
	args sortFuncArgs
}

func TestSelectionSort(t *testing.T) {
	testSortCases := []testSortCase{
		{
			name: "[1, 39, 2, 9, 7, 54, 11]",
			args: sortArgs{
				n: []int{1, 39, 2, 9, 7, 54, 11},
				w: []int{1, 2, 7, 9, 11, 39, 54},
			},
		},
		{
			name: "[4, 3, 2, 10, 12, 1, 5, 6]",
			args: sortArgs{
				n: []int{4, 3, 2, 10, 12, 1, 5, 6},
				w: []int{1, 2, 3, 4, 5, 6, 10, 12},
			},
		},
	}
	for _, tt := range testSortCases {
		t.Run(tt.name, func(t *testing.T) {
			SelectionSort(tt.args.n)
		})
		assert.Equal(t, tt.args.n, tt.args.w)
	}
}

func TestSelectionSortFunc(t *testing.T) {
	testSortFuncCases := []testSortFuncCase{
		testSortFuncCase{
			name: "[1, 39, 2, 9, 7, 54, 11] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{1, 39, 2, 9, 7, 54, 11},
					w: []int{54, 39, 11, 9, 7, 2, 1},
				},
				f: sortfunc,
			},
		},
		testSortFuncCase{
			name: "[4, 3, 2, 10, 12, 1, 5, 6] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{4, 3, 2, 10, 12, 1, 5, 6},
					w: []int{12, 10, 6, 5, 4, 3, 2, 1},
				},
				f: sortfunc,
			},
		},
	}
	for _, tt := range testSortFuncCases {
		t.Run(tt.name, func(t *testing.T) {
			SelectionSortFunc(tt.args.n, tt.args.f)
		})
		assert.Equal(t, tt.args.n, tt.args.w)
	}
}

func TestInsertionSort(t *testing.T) {
	testSortCases := []testSortCase{
		{
			name: "[1, 39, 2, 9, 7, 54, 11]",
			args: sortArgs{
				n: []int{1, 39, 2, 9, 7, 54, 11},
				w: []int{1, 2, 7, 9, 11, 39, 54},
			},
		},
		{
			name: "[4, 3, 2, 10, 12, 1, 5, 6]",
			args: sortArgs{
				n: []int{4, 3, 2, 10, 12, 1, 5, 6},
				w: []int{1, 2, 3, 4, 5, 6, 10, 12},
			},
		},
	}
	for _, tt := range testSortCases {
		t.Run(tt.name, func(t *testing.T) {
			InsertionSort(tt.args.n)
		})
		assert.Equal(t, tt.args.n, tt.args.w)
	}
}

func TestInsertionSortFunc(t *testing.T) {
	testSortFuncCases := []testSortFuncCase{
		testSortFuncCase{
			name: "[1, 39, 2, 9, 7, 54, 11] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{1, 39, 2, 9, 7, 54, 11},
					w: []int{54, 39, 11, 9, 7, 2, 1},
				},
				f: sortfunc,
			},
		},
		testSortFuncCase{
			name: "[4, 3, 2, 10, 12, 1, 5, 6] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{4, 3, 2, 10, 12, 1, 5, 6},
					w: []int{12, 10, 6, 5, 4, 3, 2, 1},
				},
				f: sortfunc,
			},
		},
	}
	for _, tt := range testSortFuncCases {
		t.Run(tt.name, func(t *testing.T) {
			InsertionSortFunc(tt.args.n, tt.args.f)
		})
		assert.Equal(t, tt.args.w, tt.args.n)
	}
}

func TestCombSort(t *testing.T) {
	testSortCases := []testSortCase{
		{
			name: "[1, 39, 2, 9, 7, 54, 11]",
			args: sortArgs{
				n: []int{1, 39, 2, 9, 7, 54, 11},
				w: []int{1, 2, 7, 9, 11, 39, 54},
			},
		},
		{
			name: "[4, 3, 2, 10, 12, 1, 5, 6]",
			args: sortArgs{
				n: []int{4, 3, 2, 10, 12, 1, 5, 6},
				w: []int{1, 2, 3, 4, 5, 6, 10, 12},
			},
		},
	}
	for _, tt := range testSortCases {
		nums := tt.args.n[0:]
		t.Run(tt.name, func(t *testing.T) {
			CombSort(nums, 1.8)
		})
		assert.Equal(t, tt.args.w, nums)
	}
}

func TestCombSortFunc(t *testing.T) {
	testSortFuncCases := []testSortFuncCase{
		testSortFuncCase{
			name: "[1, 39, 2, 9, 7, 54, 11] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{1, 39, 2, 9, 7, 54, 11},
					w: []int{54, 39, 11, 9, 7, 2, 1},
				},
				f: sortfunc,
			},
		},
		testSortFuncCase{
			name: "[4, 3, 2, 10, 12, 1, 5, 6] func(a, b int) bool { return a < b }",
			args: sortFuncArgs{
				sortArgs: sortArgs{
					n: []int{4, 3, 2, 10, 12, 1, 5, 6},
					w: []int{12, 10, 6, 5, 4, 3, 2, 1},
				},
				f: sortfunc,
			},
		},
	}
	for _, tt := range testSortFuncCases {
		t.Run(tt.name, func(t *testing.T) {
			CombSortFunc(tt.args.n, 1.8, tt.args.f)
		})
		assert.Equal(t, tt.args.w, tt.args.n)
	}
}

func TestMergeSort(t *testing.T) {
	testSortCases := []testSortCase{
		{
			name: "[1, 39, 2, 9, 7, 54, 11]",
			args: sortArgs{
				n: []int{1, 39, 2, 9, 7, 54, 11},
				w: []int{1, 2, 7, 9, 11, 39, 54},
			},
		},
		{
			name: "[4, 3, 2, 10, 12, 1, 5, 6]",
			args: sortArgs{
				n: []int{4, 3, 2, 10, 12, 1, 5, 6},
				w: []int{1, 2, 3, 4, 5, 6, 10, 12},
			},
		},
	}
	for _, tt := range testSortCases {
		t.Run(tt.name, func(t *testing.T) {
			n := MergeSort(tt.args.n)
			assert.Equal(t, tt.args.w, n)
		})
	}
}

// func TestBruteForce(t *testing.T) {
// 	type args struct {
// 		nums []int
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want [][]int
// 	}{
// 		{
// 			name: "BruteForce {1, 39, 2, 9, 7, 54, 11}",
// 			args: args{nums: []int{1, 39, 2, 9, 7, 54, 11}},
// 			want: [][]int{
// 				[]int{1, 39, 2, 9, 7, 54, 11},
// 				[]int{1, 2, 7, 9, 11, 39, 54},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := BruteForce(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("BruteForce() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
