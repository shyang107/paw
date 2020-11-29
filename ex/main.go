package main

import (
	"github.com/sirupsen/logrus"

	"github.com/shyang107/paw"
	// "github.com/thoas/go-funk"
)

var (
	// lg = paw.Glog
	lg  = paw.Logger
	log = paw.Logger
)

func init() {
	lg.SetLevel(logrus.DebugLevel)
}

func main() {
	// exLineCount()
	// exFileLineCount()
	// rehttp()
	// exGetAbbrString()
	// exTableFormat()
	// exStringBuilder()
	// exLoger()
	// exReverse()
	// exPrintTree1()
	exPrintTree2()
	// exShuffle()
	// exGetCurrPath()
	// var n1 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n2 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n3 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n4 = []int{1, 39, 2, 9, 7, 54, 11}
	// // var n1 = []int{4, 3, 2, 10, 12, 1, 5, 6}
	// // var n2 = []int{4, 3, 2, 10, 12, 1, 5, 6}
	// // size := 20
	// // n1 = paw.GenerateSlice(size)
	// InsertionSort(n1)
	// // n2 = paw.GenerateSlice(size)
	// SelectionSort(n2)
	// // n3 = paw.GenerateSlice(size)
	// exCombSort(n3)
	// // n4 = paw.GenerateSlice(size)
	// exMergeSort(n4)
	// exRegEx()
	// exLogger()
	// exFolder()
	// exGetFiles1()
	// exGetFiles2()
	// exGetFiles3()
	// exGetFilesString()
	// exGrouppingFiles1()
	// exGrouppingFiles2()
	// exGrouppingFiles3()
	exGrouppingFiles4()
	// exTextTemplate()
	// exRegEx2()
}
