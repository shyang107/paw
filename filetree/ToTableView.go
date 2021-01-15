package filetree

import (
	"fmt"

	"github.com/shyang107/paw"
)

// ToTableViewBytes will return the []byte of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToTableViewBytes(pad string) []byte {
	return []byte(f.ToTableView(pad, false))
}

// ToTableExtendViewBytes will return the []byte involving extend attribute of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToTableExtendViewBytes(pad string) []byte {
	return []byte(f.ToTableView(pad, true))
}

// ToTableView will return the string of FileList in table form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
// 	If `isExtended` is true to involve extend attribute
func (f *FileList) ToTableView(pad string, isExtended bool) string {
	var (
		// w      = new(bytes.Buffer)
		buf         = f.StringBuilder() //f.Buffer()
		w           = f.writer          //f.Writer()
		nDirs       = f.NDirs()
		nFiles      = f.NFiles()
		dirs        = f.dirs  //f.Dirs()
		fm          = f.store //f.Map()
		widthOfName = 75
	)
	buf.Reset()

	f.DisableColor()

	tf := &paw.TableFormat{
		Fields:    []string{"No.", "Mode", "Size", "Files"},
		LenFields: []int{5, 11, 6, widthOfName},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
		Padding:   pad,
		IsWrapped: true,
	}
	tf.Prepare(w)
	// tf.SetWrapFields()

	sdsize := ByteSize(f.totalSize)
	head := fmt.Sprintf("Root directory: %v, size ≈ %v", f.root, sdsize)
	tf.SetBeforeMessage(head)

	tf.PrintSart()
	j := 0
	ndirs, nfiles := 1, 0
	for i, dir := range dirs {
		sumsize := uint64(0)
		for jj, file := range fm[dir] {
			fsize := file.Size
			sfsize := ByteSize(fsize)
			mode := file.ColorPermission()
			if jj == 0 && file.IsDir() {
				idx := fmt.Sprintf("G%d", i)
				sfsize = "-"
				switch f.depth {
				case 0:
					if len(fm[dir]) > 1 && !paw.EqualFold(file.Dir, RootMark) {
						tf.PrintRow(idx, mode, sfsize, file)
					}
				default:
					name := ""
					if paw.EqualFold(file.Dir, RootMark) {
						name = file.ColorDirName("")
					} else {
						name = file.ColorDirName(f.root)
					}
					tf.PrintRow(idx, mode, sfsize, name)
				}
				continue
			}
			jdx := fmt.Sprintf("d%d", ndirs)
			name := file.Name()
			if file.IsDir() {
				ndirs++
				tf.PrintRow(jdx, mode, "-", name)
			} else {
				sumsize += fsize
				j++
				nfiles++
				tf.PrintRow(j, mode, sfsize, name)
			}
			if isExtended {
				nx := len(file.XAttributes)
				if nx > 0 {
					edge := EdgeTypeMid
					for i, x := range file.XAttributes {
						if i == nx-1 {
							edge = EdgeTypeEnd
						}
						// tf.PrintRow("", "", "", "▶︎ "+x)
						tf.PrintRow("", "", "", string(edge)+" "+x)
					}
				}
			}

		}
		if f.depth != 0 {
			// printDirSummary(buf, pad, ndirs, nfiles, sumsize)
			tf.PrintRow("", "", "", fmt.Sprintf("%v directories, %v files, size: %v.", ndirs, nfiles, ByteSize(sumsize)))

			if i != len(dirs)-1 {
				tf.PrintMiddleSepLine()
			}
		}
	}

	tf.SetAfterMessage(fmt.Sprintf("\nAccumulated %v directories, %v files, total %v.", nDirs, nFiles, ByteSize(f.totalSize)))
	tf.PrintEnd()

	f.EnableColor()

	return buf.String()
}
