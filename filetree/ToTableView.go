package filetree

import (
	"fmt"

	"github.com/thoas/go-funk"

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
		xsymb       = " @ "
		xsymb2      = "-@-"
		deftf       = &paw.TableFormat{
			Fields:            []string{"No.", "Mode", "Size", "Files"},
			LenFields:         []int{5, 11, 6, widthOfName},
			Aligns:            []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
			Padding:           pad,
			IsWrapped:         true,
			IsColorful:        true,
			XAttributeSymbol:  xsymb,
			XAttributeSymbol2: xsymb2,
		}

		nfds = len(pfields) + 1
		fds  = NewFieldSliceFrom(pfieldKeys, git)
		wNo  = paw.MaxInt(len(fmt.Sprint(nDirs))+1, len(fmt.Sprint(nFiles))+1)
		fdNo = &Field{
			Key:   PFieldNone,
			Name:  "No",
			Width: wNo,
			Align: paw.AlignRight,
			// HeadColor:  chdp,
			ValueColor: cdashp,
		}
	)
	buf.Reset()

	fds.Insert(0, fdNo)

	tf := &paw.TableFormat{
		Fields:            fds.Heads(),
		LenFields:         fds.HeadWidths(),
		Aligns:            fds.Aligns(),
		Padding:           pad,
		IsWrapped:         true,
		IsColorful:        true,
		XAttributeSymbol:  xsymb,
		XAttributeSymbol2: xsymb2,
	}

	modifyFDSWidth(fds, f, sttyWidth-2-paw.StringWidth(pad))
	fdName := fds.Get(PFieldName)

	wdmeta := fds.MetaHeadsStringWidth()
	if wdmeta > sttywd-10 {
		paw.Error.Println("too many fields")
		tf = deftf
	}

	tf.LenFields = fds.HeadWidths()

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
			fds.SetValues(file, git)
			// fds.SetColorfulValues(file, git)
			fdName.SetColorfulValue("")
			fdName.SetValueColor(GetFileLSColor(file))
			tf.Colors = fds.Colors()
			tf.FieldsColorString = fds.ColorValueStrings()

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

			if jj > 0 {
				fdNo.SetValue(jdx)
			} else { //jj==0
				fdName.SetValue(file.Dir)
			}

			values := fds.Values()
			tf.PrintRow(values...)

			if isExtended && len(file.XAttributes) > 0 {
				var (
					values = make([]interface{}, nfds)
					// nx     = len(file.XAttributes)
					wds = paw.StringWidth(xsymb)
					wdx = fdName.Width - wds //sttywd - wdmeta - wds
				)
				funk.Fill(values, "")
				for _, x := range file.XAttributes {
					wx := paw.StringWidth(x)
					if wx <= wdx {
						values[nfds-1] = tf.XAttributeSymbol + x
						tf.PrintRow(values...)
					} else {
						xs := paw.WrapToSlice(x, wdx)
						values[nfds-1] = tf.XAttributeSymbol + xs[0]
						tf.PrintRow(values...)
						for i := 1; i < len(xs); i++ {
							values[nfds-1] = tf.XAttributeSymbol2 + xs[i]
							tf.PrintRow(values...)
						}
					}
				}
			}
		}
		if f.depth != 0 {
			tf.PrintLineln(fmt.Sprintf("%s%v directories, %v files, size: %v.", pad+paw.Spaces(fdNo.Width+1), nsubdir, nsubfiles, ByteSize(sumsize)))
			if i < len(dirs)-1 && ndirs < nDirs {
				tf.PrintMiddleSepLine()
			}
		}
	}

	tf.SetAfterMessage(fmt.Sprintf("Accumulated %v directories, %v files, total %v.", nDirs, nFiles, ByteSize(f.totalSize)))
	tf.PrintEnd()

	return buf.String()
}
