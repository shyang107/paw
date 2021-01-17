package filetree

import (
	"errors"
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
		field1 = "  No."
		nfds   = len(pfields) + 1
	)
	buf.Reset()

	f.DisableColor()

	tf, err := newTable(pad, sttywd, field1, git)
	if err != nil {
		tf = deftf
	}
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
		sumsize := uint64(0)
		idx := fmt.Sprintf("G%d", i)
		if len(fm[dir]) < 2 {
			continue
		}
		for jj, file := range fm[dir] {
			jdx := ""
			if file.IsDir() {
				if !paw.EqualFold(file.Dir, RootMark) {
					ndirs++
				}
				jdx = fmt.Sprintf("d%d", ndirs)
			} else {
				sumsize += file.Size
				j++
				nfiles++
				jdx = fmt.Sprintf("%d", nfiles)
			}
			values := getRowValues(f, jj, len(fm[dir]), idx, jdx, file, git, nfds)
			// spew.Dump(values)
			tf.PrintRow(values...)
			if isExtended {
				values := make([]interface{}, nfds)
				nx := len(file.XAttributes)
				if nx > 0 {
					edge := EdgeTypeMid
					for i, x := range file.XAttributes {
						if i == nx-1 {
							edge = EdgeTypeEnd
						}
						values[nfds-1] = string(edge) + " " + x
						// values[nfds-1] = string(edge) + "▶︎ " + x
						tf.PrintRow(values...)
					}
				}
			}
		}
		if f.depth != 0 {
			// printDirSummary(buf, pad, ndirs, nfiles, sumsize)
			// tf.PrintRow("", "", "", fmt.Sprintf("%v directories, %v files, size: %v.", ndirs, nfiles, ByteSize(sumsize)))
			fmt.Fprintln(w, fmt.Sprintf("%s%v directories, %v files, size: %v.", pad+paw.Spaces(paw.StringWidth(field1)+1), ndirs, nfiles, ByteSize(sumsize)))
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
func getRowValues(fl *FileList, j, nfiles int, idx, jdx string, file *File, git GitStatus, nfds int) (values []interface{}) {
	sfsize := "-"
	if j == 0 && file.IsDir() {
		values = append(values, idx)
		for k := 0; k < len(pfields); k++ {
			key := pfieldKeys[k]
			switch key {
			case PFieldGit:
				if git.NoGit {
					continue
				} else {
					values = append(values, getFileValue(file, key, git))
				}
			case PFieldSize:
				if fl.depth == 0 && nfiles > 1 && !paw.EqualFold(file.Dir, RootMark) {
					if key&PFieldSize != 0 {
						values = append(values, sfsize)
					}
				} else {
					values = append(values, getFileValue(file, key, git))
				}
			case PFieldName:
				dir, name := "", ""
				if paw.EqualFold(file.Dir, RootMark) {
					dir, name = getDirAndName(file.Path, "")
				} else {
					dir, name = getDirAndName(file.Path, fl.root)
				}
				if file.IsLink() {
					values = append(values, dir+" -> "+name)
				} else {
					values = append(values, dir+name)
				}
			default:
				values = append(values, getFileValue(file, key, git))
			}
		}
		return values
	}

	values = append(values, jdx)
	for k := 0; k < len(pfields); k++ {
		key := pfieldKeys[k]
		switch key {
		case PFieldGit:
			if git.NoGit {
				continue
			} else {
				values = append(values, getFileValue(file, key, git))
			}
		case PFieldSize:
			if file.IsDir() {
				values = append(values, sfsize)
			} else {
				values = append(values, getFileValue(file, key, git))
			}
		default:
			values = append(values, getFileValue(file, key, git))
		}
	}
	return values
}

func getFileValue(file *File, key PDFieldFlag, git GitStatus) (value interface{}) {
	switch key {
	case PFieldINode: //"inode",
		value = file.INode()
	case PFieldPermissions: //"Permissions",
		value = file.Permission()
	case PFieldLinks: //"Links",
		value = file.NLinks()
	case PFieldSize: //"Size",
		value = ByteSize(file.Size)
	case PFieldBlocks: //"Blocks",
		value = file.NLinks()
	case PFieldUser: //"User",
		value = urname
	case PFieldGroup: //"Group",
		value = gpname
	case PFieldModified: //"Date Modified",
		value = DateString(file.ModifiedTime())
	case PFieldCreated: //"Date Created",
		value = DateString(file.CreatedTime())
	case PFieldAccessed: //"Date Accessed",
		value = DateString(file.AccessedTime())
	case PFieldGit: //"Git",
		value = getGitStatus(git, file)
	case PFieldName: //"Name",
		value = file.Name()
	}
	return value
}

func newTable(pad string, sttywd int, field1 string, git GitStatus) (*paw.TableFormat, error) {
	var (
		nfds = len(pfields) + 1
	)
	var fds []string
	var lenfds []int
	var aligns []paw.Align
	fds = append(fds, field1)
	lenfds = append(lenfds, len(field1))
	aligns = append(aligns, paw.AlignRight)
	for _, k := range pfieldKeys {
		switch k {
		case PFieldPermissions, PFieldUser, PFieldGroup, PFieldModified, PFieldCreated, PFieldAccessed, PFieldName:
			fds = append(fds, pfieldsMap[k])
			lenfds = append(lenfds, pfieldWidthsMap[k])
			aligns = append(aligns, paw.AlignLeft)
		case PFieldGit:
			if git.NoGit {
				continue
			} else {
				fds = append(fds, pfieldsMap[k])
				lenfds = append(lenfds, pfieldWidthsMap[k])
				aligns = append(aligns, paw.AlignRight)
			}
		default:
			fds = append(fds, pfieldsMap[k])
			lenfds = append(lenfds, pfieldWidthsMap[k])
			aligns = append(aligns, paw.AlignRight)
		}
	}
	nfds = len(lenfds)
	lenfds[nfds-1] = sttywd - paw.SumInts(lenfds[:nfds-1]...) - nfds + 1

	if paw.SumInts(lenfds...) > sttywd {
		return nil, errors.New("too many fields")
	}
	tf := &paw.TableFormat{
		Fields:    fds,
		LenFields: lenfds,
		Aligns:    aligns,
		Padding:   pad,
		IsWrapped: true,
	}
	return tf, nil
}
