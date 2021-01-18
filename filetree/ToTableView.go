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
		git         = f.GetGitStatus()
		dirs        = f.dirs  //f.Dirs()
		fm          = f.store //f.Map()
		sttywd      = sttyWidth - 2
		widthOfName = 75
		deftf       = &paw.TableFormat{
			Fields:    []string{"No.", "Mode", "Size", "Files"},
			LenFields: []int{5, 11, 6, widthOfName},
			Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
			Padding:   pad,
			IsWrapped: true,
		}

		nfds = len(pfields) + 1
		fds  = NewFieldSliceFrom(pfieldKeys, git)
		wNo  = len(fmt.Sprint(f.NFiles())) + 1
		fdNo = &Field{
			Key:   PFieldNone,
			Name:  "No",
			Width: wNo,
			Align: paw.AlignRight,
			// headcp: chdp,
		}
	)
	buf.Reset()
	// f.DisableColor()

	fds.Insert(0, fdNo)

	tf := &paw.TableFormat{
		Fields:    fds.Heads(),
		LenFields: fds.HeadWidths(),
		Aligns:    fds.HeadAligns(),
		Padding:   pad,
		IsWrapped: true,
	}

	wdmeta := fds.MetaValuesStringWidth()
	if wdmeta > sttywd-10 {
		tf = deftf
	} else {
		fds.Get(PFieldName).Width = sttywd - wdmeta
	}
	modifySizeWidth(fds, f)

	// spmeta := pad + paw.Spaces(wdmeta)

	tf.LenFields = fds.HeadWidths()

	fdName := fds.Get(PFieldName)

	nfds = tf.NFields()
	tf.Prepare(w)
	// tf.SetWrapFields()

	sdsize := ByteSize(f.totalSize)
	head := fmt.Sprintf("Root directory: %v, size ≈ %v", f.root, sdsize)
	tf.SetBeforeMessage(head)

	tf.PrintSart()
	j := 0
	ndirs, nfiles := 0, 0
	for i, dir := range dirs {
		idx := fmt.Sprintf("G%d", i)
		fdNo.SetValue(idx)
		if len(fm[dir]) < 2 {
			continue
		}
		nsubdir, nsubfiles, sumsize := 0, 0, uint64(0)
		for jj, file := range fm[dir] {
			jdx := ""
			if file.IsDir() {
				if jj != 0 && !paw.EqualFold(file.Dir, RootMark) {
					ndirs++
					nsubdir++
				}
				jdx = fmt.Sprintf("d%d", ndirs)
			} else {
				sumsize += file.Size
				j++
				nfiles++
				nsubfiles++
				jdx = fmt.Sprintf("%d", nfiles)
			}

			fds.SetValues(file, git)

			if jj > 0 {
				fdNo.SetValue(jdx)
			} else { //jj==0
				fdName.SetValue(file.Dir)
			}

			values := fds.Values()
			tf.PrintRow(values...)

			if isExtended && len(file.XAttributes) > 0 {
				values := make([]interface{}, nfds)
				nx := len(file.XAttributes)
				if nx > 0 {
					edge := EdgeTypeMid
					for i, x := range file.XAttributes {
						if i == nx-1 {
							edge = EdgeTypeEnd
						}
						wde := edgeWidth[edge] + 1
						wdx := sttywd - wdmeta - wde - 3
						wx := paw.StringWidth(x)
						if wx <= wdx {
							values[nfds-1] = string(edge) + " " + x
							tf.PrintRow(values...)
						} else {
							xs := paw.Split(paw.Wrap(x, wdx), "\n")
							values[nfds-1] = string(edge) + " " + xs[0]
							tf.PrintRow(values...)
							padx := ""
							switch edge {
							case EdgeTypeMid:
								padx = fmt.Sprintf("%s%s", EdgeTypeLink, SpaceIndentSize)
							case EdgeTypeEnd:
								padx = fmt.Sprintf("%s ", paw.Spaces(edgeWidth[edge]))
							}
							for i := 1; i < len(xs); i++ {
								values[nfds-1] = padx + xs[i]
								tf.PrintRow(values...)
							}

						}
					}
				}
			}
		}
		if f.depth != 0 {
			fmt.Fprintln(w, fmt.Sprintf("%s%v directories, %v files, size: %v.", pad+paw.Spaces(fdNo.Width+1), nsubdir, nsubfiles, ByteSize(sumsize)))
			if i < len(dirs)-1 && ndirs < nDirs {
				tf.PrintMiddleSepLine()
			}
		}
	}

	tf.SetAfterMessage(fmt.Sprintf("Accumulated %v directories, %v files, total %v.", nDirs, nFiles, ByteSize(f.totalSize)))
	tf.PrintEnd()

	// f.EnableColor()

	return buf.String()
}
