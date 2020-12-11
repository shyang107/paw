package _junk

import "sort"

// LessFunc implement Less()
type LessFunc func(p1, p2 *File) bool

// FilesSorter implements the Sort interface, sorting the files within.
type FilesSorter struct {
	files []File
	less  []LessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *FilesSorter) Sort(files []File) {
	ms.files = files
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...LessFunc) *FilesSorter {
	return &FilesSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *FilesSorter) Len() int {
	return len(ms.files)
}

// Swap is part of sort.Interface.
func (ms *FilesSorter) Swap(i, j int) {
	ms.files[i], ms.files[j] = ms.files[j], ms.files[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *FilesSorter) Less(i, j int) bool {
	p, q := &ms.files[i], &ms.files[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
