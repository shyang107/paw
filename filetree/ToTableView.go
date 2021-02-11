package filetree

import (
	"fmt"
	"strings"

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
		w                = f.StringBuilder() //f.Buffer()
		nDirs, nFiles, _ = f.NTotalDirsAndFile()
		nItems           = nDirs + nFiles
		wdidx            = len(fmt.Sprint(nDirs))
		wdjdx            = paw.MaxInts(wdidx, len(fmt.Sprint(nFiles)))
		git              = f.GetGitStatus()
		dirs             = f.dirs  //f.Dirs()
		fm               = f.store //f.Map()
		wpad             = paw.StringWidth(pad)
		wdstty           = sttyWidth - 2 - wpad
		banner           = strings.Repeat("-", wdstty)
		widthOfName      = 75
		// xsymb            = paw.XAttrSymbol
		// xsymb2           = paw.XAttrSymbol2
		deftf = &paw.TableFormat{
			Fields:            []string{"No.", "Mode", "Size", "Files"},
			LenFields:         []int{5, 11, 6, widthOfName},
			Aligns:            []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
			Padding:           pad,
			IsWrapped:         true,
			IsColorful:        true,
			XAttributeSymbol:  paw.XAttrSymbol,
			XAttributeSymbol2: paw.XAttrSymbol2,
		}

		// nfds = len(pfields) + 1
		fds  = NewFieldSliceFrom(pfieldKeys, git)
		wNo  = paw.MaxInt(wdidx, wdjdx) + 1
		fdNo = &Field{
			Key:   PFieldNone,
			Name:  "No",
			Width: wNo,
			Align: paw.AlignRight,
			// HeadColor:  chdp,
			ValueColor: cdashp,
		}

		spNo     = paw.Spaces(fdNo.Width + 1)
		roothead = getColorizedRootHead(f.root, f.TotalSize(), wdstty)
	)

	w.Reset()

	fds.Insert(0, fdNo)
	fds.ModifyWidth(f, wdstty)
	fdName := fds.Get(PFieldName)

	tf := &paw.TableFormat{
		Fields:            fds.Heads(false),
		LenFields:         fds.Widths(),
		Aligns:            fds.Aligns(),
		Padding:           pad,
		IsWrapped:         true,
		IsColorful:        true,
		XAttributeSymbol:  paw.XAttrSymbol,
		XAttributeSymbol2: paw.XAttrSymbol2,
	}

	wdmeta := fds.MetaHeadsStringWidth()
	if wdmeta > wdstty-10 {
		paw.Warning.Println("too many fields, use default table view")
		tf = deftf
	}

	tf.LenFields = fds.Widths()

	// nfds = tf.NFields()
	tf.Prepare(w)
	// tf.SetWrapFields()

	tf.SetBeforeMessage(roothead)

	tf.PrintSart()

	var ndirs, nfiles = 0, 0
	for i, dir := range dirs {
		idx := fmt.Sprintf("G%-[1]*[2]d ", wdidx, i)
		widx := paw.StringWidth(idx)
		cidx := cdip.Sprint(idx)
		if len(fm[dir]) > 1 {
			if !strings.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					tf.PrintLine(fm[dir][0].DirNameWrapC(cidx, wdstty-widx))
					errmsg := f.GetErrorString(dir)
					if len(errmsg) > 0 {
						tf.PrintLine(errmsg)
					}
				}
			}
		} else {
			continue
		}
		if len(fm[dir]) > 1 &&
			f.depth != 0 &&
			!strings.EqualFold(dir, RootMark) {
			tf.Fields = fds.Heads(false)
			tf.PrintHeads()
		}

		for _, file := range fm[dir][1:] {
			cjdx, jdx := "", ""
			if file.IsDir() {
				ndirs++
				jdx = fmt.Sprintf("D%d", ndirs)
				cjdx = cdip.Sprintf("%[1]*[2]s", wNo, jdx)
			} else {
				nfiles++
				jdx = fmt.Sprintf("F%d", nfiles)
				cjdx = cfip.Sprintf("%[1]*[2]s", wNo, jdx)
			}
			fds.SetValues(file, git)

			fdNo.SetValueC(cjdx)
			fdNo.Width = wNo

			fdName.SetValueColor(GetFileLSColor(file))
			fdName.SetValue(file.Name())

			tf.Colors = fds.Colors()
			tf.FieldsColorString = fds.ValueStringCs()
			values := fds.Values()
			tf.PrintRow(values...)

			if isExtended && len(file.XAttributes) > 0 {
				var (
					wds = paw.StringWidth(tf.XAttributeSymbol)
					wdx = fdName.Width - wds
				)
				for _, x := range file.XAttributes {
					fds.EmptyValues()
					wx := paw.StringWidth(x)
					if wx <= wdx {
						fdName.Value = tf.XAttributeSymbol + x
						values := fds.Values()
						tf.PrintRow(values...)
					} else {
						xs := paw.WrapToSlice(x, wdx)
						fdName.Value = tf.XAttributeSymbol + xs[0]
						values := fds.Values()
						tf.PrintRow(values...)
						for i := 1; i < len(xs); i++ {
							fdName.Value = tf.XAttributeSymbol2 + xs[i]
							values := fds.Values()
							tf.PrintRow(values...)
						}
					}
				}
			}
		}
		if f.depth != 0 {
			tf.PrintLineln(cpmpt.Sprint(pad+spNo) + f.DirSummary(dir, wdstty-paw.StringWidth(pad+spNo)))

			if ndirs+nfiles < nItems {
				tf.PrintLineln(banner)
			}
		}
	}

	tf.SetAfterMessage(f.TotalSummary(wdstty))
	tf.PrintEnd()

	return w.String()
}
