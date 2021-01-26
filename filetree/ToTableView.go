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
		buf              = f.StringBuilder() //f.Buffer()
		w                = f.writer          //f.Writer()
		nDirs, nFiles, _ = f.NTotalDirsAndFile()
		nItems           = nDirs + nFiles
		wdidx            = len(fmt.Sprint(nDirs))
		wdjdx            = paw.MaxInts(wdidx, len(fmt.Sprint(nFiles)))
		git              = f.GetGitStatus()
		dirs             = f.dirs  //f.Dirs()
		fm               = f.store //f.Map()
		wpad             = paw.StringWidth(pad)
		wstty            = sttyWidth - 2 - wpad
		banner           = paw.Repeat("-", wstty)
		widthOfName      = 75
		xsymb            = paw.XAttrSymbol
		xsymb2           = paw.XAttrSymbol2
		deftf            = &paw.TableFormat{
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
		rootName = GetColorizedDirName(f.root, "")
		ctdsize  = GetColorizedSize(f.totalSize)
		head     = fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, rootName, ctdsize)
	)
	buf.Reset()

	fds.Insert(0, fdNo)

	tf := &paw.TableFormat{
		Fields:            fds.Heads(false),
		LenFields:         fds.Widths(),
		Aligns:            fds.Aligns(),
		Padding:           pad,
		IsWrapped:         true,
		IsColorful:        true,
		XAttributeSymbol:  xsymb,
		XAttributeSymbol2: xsymb2,
	}

	fds.ModifyWidth(f, wstty)

	fdName := fds.Get(PFieldName)

	wdmeta := fds.MetaHeadsStringWidth()
	if wdmeta > wstty-10 {
		paw.Warning.Println("too many fields, use default table view")
		tf = deftf
	}

	tf.LenFields = fds.Widths()

	nfds = tf.NFields()
	tf.Prepare(w)
	// tf.SetWrapFields()

	tf.SetBeforeMessage(head)

	tf.PrintSart()

	var ndirs, nfiles = 0, 0
	for i, dir := range dirs {
		idx := fmt.Sprintf("G%-[1]*[2]d ", wdidx, i)
		widx := paw.StringWidth(idx)
		cidx := cdip.Sprint(idx)
		if len(fm[dir]) > 1 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					tf.PrintLine(fm[dir][0].ColorWrapDirName(cidx, wstty-widx))
				}
			}
		} else {
			continue
		}
		if len(fm[dir]) > 1 &&
			f.depth != 0 &&
			!paw.EqualFold(dir, RootMark) {
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

			fdNo.SetColorfulValue(cjdx)
			fdNo.Width = wNo

			fdName.SetValueColor(GetFileLSColor(file))
			fdName.SetValue(file.Name())

			tf.Colors = fds.Colors()
			tf.FieldsColorString = fds.ColorValueStrings()
			values := fds.Values()
			tf.PrintRow(values...)

			if isExtended && len(file.XAttributes) > 0 {
				var (
					values = make([]interface{}, nfds)
					wds    = paw.StringWidth(xsymb)
					wdx    = fdName.Width - wds
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
			tf.PrintLineln(pad + spNo + f.DirSummary(dir))

			if ndirs+nfiles < nItems {
				tf.PrintLineln(banner)
			}
		}
	}

	tf.SetAfterMessage(f.TotalSummary())
	tf.PrintEnd()

	return buf.String()
}
