package vfs

// // PrintDirOption is the option of PrintDir
// //
// // Fields:
// // 	Depth:
// // 		Depth < 0 : print all files and directories recursively of argument path of PrintDir.
// // 		Depth = 0 : print files and directories only in argument path of PrintDir.
// // 		Depth > 0 : print files and directories recursively under depth of directory in argument path of PrintDir.
// // ViewFlag: the view-option of PrintDir
// // Call
// type PrintDirOption struct {
// 	Depth     int
// 	ViewFlag  PDViewFlag
// 	FieldFlag PDFieldFlag
// 	SortOpt   *PDSortOption
// 	FiltOpt   *PDFilterOption
// 	Ignore    IgnoreFunc
// 	//
// 	Root      string
// 	Paths     []string
// 	isTrace   bool
// 	File      *File
// 	fieldKeys []PDFieldFlag
// 	// fields      []string
// 	// fieldWidths []int
// 	isGit bool
// }

// func (p *PrintDirOption) FieldKeys() []PDFieldFlag {
// 	return p.fieldKeys
// }

// func (p *PrintDirOption) Fields() []string {
// 	names := []string{}
// 	for _, v := range p.fieldKeys {
// 		names = append(names, v.Name())
// 	}
// 	return names
// }

// func (p *PrintDirOption) FieldWidths() []int {
// 	wds := []int{}
// 	for _, v := range p.fieldKeys {
// 		wds = append(wds, len(v.Name()))
// 	}
// 	return wds
// }

// func (p *PrintDirOption) ConfigFields() {
// 	paw.Logger.Info()

// 	p.fieldKeys = []PDFieldFlag{}
// 	// p.fields = []string{}
// 	// p.fieldWidths = []int{}
// 	if p.FieldFlag&PFieldINode != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldINode)
// 	}

// 	p.fieldKeys = append(p.fieldKeys, PFieldPermissions)

// 	if p.FieldFlag&PFieldLinks != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldLinks)
// 	}

// 	p.fieldKeys = append(p.fieldKeys, PFieldSize)

// 	if p.FieldFlag&PFieldBlocks != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldBlocks)
// 	}

// 	p.fieldKeys = append(p.fieldKeys, PFieldUser)
// 	p.fieldKeys = append(p.fieldKeys, PFieldGroup)

// 	if p.FieldFlag&PFieldModified != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldModified)
// 	}
// 	if p.FieldFlag&PFieldCreated != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldCreated)
// 	}
// 	if p.FieldFlag&PFieldAccessed != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldAccessed)
// 	}

// 	if p.FieldFlag&PFieldMd5 != 0 {
// 		hasMd5 = true
// 		p.fieldKeys = append(p.fieldKeys, PFieldMd5)
// 	}

// 	if p.FieldFlag&PFieldGit != 0 {
// 		p.fieldKeys = append(p.fieldKeys, PFieldGit)
// 		p.isGit = true
// 	}
// 	// p.fieldKeys = append(p.fieldKeys, PFieldGit)

// 	p.fieldKeys = append(p.fieldKeys, PFieldName)

// 	// for _, k := range p.fieldKeys {
// 	// 	p.fields = append(p.fields, k.Name())
// 	// 	p.fieldWidths = append(p.fieldWidths, k.Width())
// 	// }
// }

// func (p *PrintDirOption) SetDepth(depth int) {
// 	p.Depth = depth
// }

// func (p *PrintDirOption) SetViewFlag(viewFlag PDViewFlag) {
// 	p.ViewFlag = viewFlag
// }

// func (p *PrintDirOption) SetFieldFlag(fieldFlag PDFieldFlag) {
// 	p.FieldFlag = fieldFlag
// }

// func (p *PrintDirOption) SetSortOpt(sortOption *PDSortOption) {
// 	p.SortOpt = sortOption
// }

// func (p *PrintDirOption) SetFiltOpt(filterOption *PDFilterOption) {
// 	p.FiltOpt = filterOption
// }

// func (p *PrintDirOption) SetIgnore(f IgnoreFunc) {
// 	p.Ignore = f
// }

// func (p *PrintDirOption) ConfigFilter() {
// 	paw.Logger.Info()

// 	igfunc := p.Ignore
// 	filtOpt := p.FiltOpt
// 	if filtOpt != nil && filtOpt.IsFilt {
// 		switch filtOpt.FiltWay {
// 		case PDFiltNoEmptyDir: // no empty dir
// 			p.Ignore = func(f *File, err error) error {
// 				if errig := igfunc(f, err); errig != nil {
// 					return errig
// 				}
// 				fis, errfilt := os.ReadDir(f.Path)
// 				if errfilt != nil {
// 					return errfilt
// 				}
// 				if len(fis) == 0 {
// 					return SkipThis
// 				}
// 				if f.IsDir() {
// 					nfiles := 0
// 					filepath.WalkDir(f.Path, func(path string, d fs.DirEntry, err error) error {
// 						if !d.IsDir() {
// 							nfiles++
// 						}
// 						return nil
// 					})
// 					if nfiles == 0 {
// 						return SkipThis
// 					}
// 				}
// 				return nil
// 			}
// 		case PDFiltJustDirs: // no files
// 			p.Ignore = func(f *File, err error) error {
// 				if errig := igfunc(f, err); errig != nil {
// 					return errig
// 				}
// 				if !f.IsDir() {
// 					return SkipThis
// 				}
// 				return nil
// 			}
// 		case PDFiltJustFiles, PDFiltJustFilesButNoEmptyDir: // PDFiltJustFilesButNoEmptyDir // no dirs
// 			p.Ignore = func(f *File, err error) error {
// 				if errig := igfunc(f, err); errig != nil {
// 					return errig
// 				}
// 				if f.IsDir() {
// 					if p.Depth == 0 {
// 						return SkipThis
// 					}
// 					nfiles := 0
// 					filepath.WalkDir(f.Path, func(path string, d fs.DirEntry, err error) error {
// 						idepth := len(strings.Split(strings.Replace(path, p.Root, ".", 1), PathSeparator)) - 1
// 						if p.Depth > 0 {
// 							if idepth > p.Depth {
// 								return SkipThis
// 							}
// 						}
// 						if !d.IsDir() {
// 							nfiles++
// 						}
// 						return nil
// 					})
// 					if nfiles == 0 {
// 						return SkipThis
// 					}
// 				}
// 				return nil
// 			}
// 		case PDFiltJustDirsButNoEmpty: // no file and no empty dir
// 			p.Ignore = func(f *File, err error) error {
// 				if errig := igfunc(f, nil); errig != nil {
// 					return errig
// 				}
// 				if f.IsDir() {
// 					nfiles := 0
// 					filepath.WalkDir(f.Path, func(path string, d fs.DirEntry, err error) error {
// 						if !d.IsDir() {
// 							nfiles++
// 						}
// 						return nil
// 					})
// 					if nfiles == 0 {
// 						return SkipThis
// 					}
// 					return nil
// 				} else {
// 					return SkipThis
// 				}
// 			}
// 		}
// 	}
// 	p.FiltOpt.IsFilt = false
// }

// func (p *PrintDirOption) SetRoot(path string) {
// 	p.Root = path
// }

// func (p *PrintDirOption) AddPath(path string) {
// 	if p.Paths == nil {
// 		p.Paths = []string{}
// 	}
// 	p.Paths = append(p.Paths, path)
// }

// func (p *PrintDirOption) NPath() int {
// 	if p.Paths == nil {
// 		return -1
// 	}
// 	return len(p.Paths)
// }

// func (p *PrintDirOption) EnableTrace(isTrace bool) {
// 	p.isTrace = isTrace
// }

// func NewPrintDirOption() *PrintDirOption {
// 	p := &PrintDirOption{
// 		Depth:     0,
// 		ViewFlag:  PListView,
// 		FieldFlag: PFieldDefault,
// 		// SortOpt:
// 		// FiltOpt:,
// 		Ignore:  DefaultIgnoreFn,
// 		Root:    "",
// 		Paths:   nil,
// 		isTrace: false,
// 	}
// 	p.ConfigFields()
// 	return p
// }

// // PDSortOption defines sorting way view of PrintDir
// //
// // Defaut:
// //  increasing sort by lower name of path
// type PDSortOption struct {
// 	IsSort   bool
// 	Reverse  bool
// 	SortFlag PDSortFlag
// }

// type PDSortFlag int

// const (
// 	PDSort PDSortFlag = 1 << iota
// 	PDSortReverse
// 	pdSortKeyINode
// 	pdSortKeyLinks
// 	pdSortKeySize
// 	pdSortKeyBlocks
// 	pdSortKeyMTime
// 	pdSortKeyCTime
// 	pdSortKeyATime
// 	pdSortKeyName

// 	PDSortByINode  = PDSort | pdSortKeyINode
// 	PDSortByINodeR = PDSortByINode | PDSortReverse

// 	PDSortByLinks  = PDSort | pdSortKeyLinks
// 	PDSortByLinksR = PDSortByLinks | PDSortReverse

// 	PDSortBySize  = PDSort | pdSortKeySize
// 	PDSortBySizeR = PDSortBySize | PDSortReverse

// 	PDSortByBlocks  = PDSort | pdSortKeyBlocks
// 	PDSortByBlocksR = PDSortByBlocks | PDSortReverse

// 	PDSortByMTime  = PDSort | pdSortKeyMTime
// 	PDSortByMTimeR = PDSortByMTime | PDSortReverse

// 	PDSortByCTime  = PDSort | pdSortKeyCTime
// 	PDSortByCTimeR = PDSortByCTime | PDSortReverse

// 	PDSortByATime  = PDSort | pdSortKeyATime
// 	PDSortByATimeR = PDSortByATime | PDSortReverse

// 	PDSortByName  = PDSort | pdSortKeyName
// 	PDSortByNameR = PDSortByName | PDSortReverse
// )

// func (s PDSortFlag) String() string {
// 	switch s {
// 	case PDSort:
// 		return "Sort"
// 	case PDSortReverse:
// 		return "Sort reversely"
// 	// case pdSortKeyINode:
// 	// 	return "Sort-Key: inode"
// 	// case pdSortKeyLinks:
// 	// 	return "Sort-Key: Links"
// 	// case pdSortKeySize:
// 	// 	return "Sort-Key: Size"
// 	// case pdSortKeyBlocks:
// 	// 	return "Sort-Key: Blocks"
// 	// case pdSortKeyMTime:
// 	// 	return "Sort-Key: Modified"
// 	// case pdSortKeyCTime:
// 	// 	return "Sort-Key: Created"
// 	// case pdSortKeyATime:
// 	// 	return "Sort-Key: Accessed"
// 	// case pdSortKeyName:
// 	// 	return "Sort-Key: Name"
// 	case PDSortByINode, pdSortKeyINode: // PDSort | pdSortKeyINode
// 		return "by inode"
// 	case PDSortByINodeR: // PDSortByINode | PDSortReverse
// 		return "by inode reversely"
// 	case PDSortByLinks, pdSortKeyLinks: // PDSort | pdSortKeyLinks
// 		return "by Links"
// 	case PDSortByLinksR: // PDSortByLinks | PDSortReverse
// 		return "by Links reversely"
// 	case PDSortBySize, pdSortKeySize: // PDSort | pdSortKeySize
// 		return "by Size"
// 	case PDSortBySizeR: // PDSortBySize | PDSortReverse
// 		return "by Size reversely"
// 	case PDSortByBlocks, pdSortKeyBlocks: // PDSort | pdSortKeyBlocks
// 		return "by Blocks"
// 	case PDSortByBlocksR: // PDSortByBlocks | PDSortReverse
// 		return "by Blocks reversely"
// 	case PDSortByMTime, pdSortKeyMTime: // PDSort | pdSortKeyMTime
// 		return "by Modified"
// 	case PDSortByMTimeR: // PDSortByMTime | PDSortReverse
// 		return "by Modified reversely"
// 	case PDSortByCTime, pdSortKeyCTime: // PDSort | pdSortKeyCTime
// 		return "by Created"
// 	case PDSortByCTimeR: // PDSortByCTime | PDSortReverse
// 		return "by Created reversely"
// 	case PDSortByATime, pdSortKeyATime: // PDSort | pdSortKeyATime
// 		return "by Accessed"
// 	case PDSortByATimeR: // PDSortByATime | PDSortReverse
// 		return "by Accessed reversely"
// 	case PDSortByName, pdSortKeyName: // PDSort | pdSortKeyName
// 		return "by Name"
// 	case PDSortByNameR: // PDSortByName | PDSortReverse
// 		return "by Name reversely"
// 	default:
// 		return "Unknown"
// 	}
// }
// func (s PDSortFlag) By() FilesBy {
// 	switch s {
// 	case PDSortReverse:
// 		return byNameR
// 	case PDSortByINode: // PDSort | pdSortKeyINode
// 		return byINode
// 	case PDSortByINodeR: // PDSortByINode | PDSortReverse
// 		return byINodeR
// 	case PDSortByLinks: // PDSort | pdSortKeyLinks
// 		return byLinks
// 	case PDSortByLinksR: // PDSortByLinks | PDSortReverse
// 		return byLinksR
// 	case PDSortBySize: // PDSort | pdSortKeySize
// 		return bySize
// 	case PDSortBySizeR: // PDSortBySize | PDSortReverse
// 		return bySizeR
// 	case PDSortByBlocks: // PDSort | pdSortKeyBlocks
// 		return byBlocks
// 	case PDSortByBlocksR: // PDSortByBlocks | PDSortReverse
// 		return byBlocksR
// 	case PDSortByMTime: // PDSort | pdSortKeyMTime
// 		return byMTime
// 	case PDSortByMTimeR: // PDSortByMTime | PDSortReverse
// 		return byMTimeR
// 	case PDSortByCTime: // PDSort | pdSortKeyCTime
// 		return byCTime
// 	case PDSortByCTimeR: // PDSortByCTime | PDSortReverse
// 		return byCTimeR
// 	case PDSortByATime: // PDSort | pdSortKeyATime
// 		return byATime
// 	case PDSortByATimeR: // PDSortByATime | PDSortReverse
// 		return byATimeR
// 	case PDSortByNameR: // PDSortByName | PDSortReverse
// 		return byNameR
// 	default: // PDSortByName: // PDSort | pdSortKeyName
// 		return byName
// 	}
// }

// func (s PDSortFlag) GetFlag(flag string) PDSortFlag {
// 	flag = strings.ToLower(flag)
// 	if flag, ok := sortMapFlag[flag]; ok {
// 		s = flag
// 	} else {
// 		s = PDSortByName
// 	}
// 	return s
// }

// var (
// 	sortMapFlag = map[string]PDSortFlag{
// 		"inode":     PDSortByINode,
// 		"links":     PDSortByLinks,
// 		"size":      PDSortBySize,
// 		"blocks":    PDSortByBlocks,
// 		"modified":  PDSortByMTime,
// 		"mtime":     PDSortByMTime,
// 		"accessed":  PDSortByATime,
// 		"atime":     PDSortByATime,
// 		"created":   PDSortByCTime,
// 		"ctime":     PDSortByCTime,
// 		"name":      PDSortByName,
// 		"inoder":    PDSortByINodeR,
// 		"linksr":    PDSortByLinksR,
// 		"sizer":     PDSortBySizeR,
// 		"blocksr":   PDSortByBlocksR,
// 		"modifiedr": PDSortByMTimeR,
// 		"mtimer":    PDSortByMTimeR,
// 		"accessedr": PDSortByATimeR,
// 		"atimer":    PDSortByATimeR,
// 		"createdr":  PDSortByCTimeR,
// 		"ctimer":    PDSortByCTimeR,
// 		"namer":     PDSortByNameR,
// 	}
// )

// // var (
// // // sortedFields = []string{"size", "modified", "accessed", "created", "name"}

// // // sortByField = map[PDSortFlag]FilesBy{
// // // 	PDSortByINode:   byINode,
// // // 	PDSortByINodeR:  byINodeR,
// // // 	PDSortByLinks:   byLinks,
// // // 	PDSortByLinksR:  byLinksR,
// // // 	PDSortBySize:    bySize,
// // // 	PDSortBySizeR:   bySizeR,
// // // 	PDSortByBlocks:  byBlocks,
// // // 	PDSortByBlocksR: byBlocksR,
// // // 	PDSortByMTime:   byMTime,
// // // 	PDSortByMTimeR:  byMTimeR,
// // // 	PDSortByATime:   byATime,
// // // 	PDSortByATimeR:  byATimeR,
// // // 	PDSortByCTime:   byCTime,
// // // 	PDSortByCTimeR:  byCTimeR,
// // // 	PDSortByName:    byName,
// // // 	PDSortByNameR:   byNameR,
// // // }
// // )

// type PDFiltFlag int

// const (
// 	PDFiltNoEmptyDir = 1 << iota
// 	PDFiltJustDirs
// 	PDFiltJustFiles
// 	PDFiltJustDirsButNoEmpty     //= PDFiltNoEmptyDir | PDFiltJustDirs
// 	PDFiltJustFilesButNoEmptyDir //= PDFiltJustFiles
// )

// func (f PDFiltFlag) String() string {
// 	switch f {
// 	case PDFiltNoEmptyDir:
// 		return "No Empty Directories"
// 	case PDFiltJustDirs:
// 		return "Just Directories"
// 	case PDFiltJustFiles:
// 		return "Just Files"
// 	case PDFiltJustDirsButNoEmpty:
// 		return "Just Directories; but no empty"
// 	case PDFiltJustFilesButNoEmptyDir:
// 		return "Just Files"
// 	default:
// 		return "Unknown"
// 	}
// }

// type PDFilterOption struct {
// 	IsFilt  bool
// 	FiltWay PDFiltFlag
// }
