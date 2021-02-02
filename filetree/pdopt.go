package filetree

// PrintDirOption is the option of PrintDir
//
// Fields:
// 	Depth:
// 		Depth < 0 : print all files and directories recursively of argument path of PrintDir.
// 		Depth = 0 : print files and directories only in argument path of PrintDir.
// 		Depth > 0 : print files and directories recursively under depth of directory in argument path of PrintDir.
// OutOpt: the view-option of PrintDir
// Call
type PrintDirOption struct {
	Depth     int
	OutOpt    PDViewFlag
	FieldFlag PDFieldFlag
	SortOpt   *PDSortOption
	FiltOpt   *PDFilterOption
	Ignore    IgnoreFunc
	//
	Root  string
	Paths []string
}

func (p *PrintDirOption) SetDepth(depth int) {
	p.Depth = depth
}

func (p *PrintDirOption) SetOutOpt(flag PDViewFlag) {
	p.OutOpt = flag
}

func (p *PrintDirOption) SetFieldFlag(flag PDFieldFlag) {
	p.FieldFlag = flag
}

func (p *PrintDirOption) SetSortOpt(opt *PDSortOption) {
	p.SortOpt = opt
}

func (p *PrintDirOption) SetFiltOpt(opt *PDFilterOption) {
	p.FiltOpt = opt
}

func (p *PrintDirOption) SetIgnore(f IgnoreFunc) {
	p.Ignore = f
}

func (p *PrintDirOption) SetRoot(path string) {
	p.Root = path
}

func (p *PrintDirOption) AddPath(path string) {
	if p.Paths == nil {
		p.Paths = []string{}
	}
	p.Paths = append(p.Paths, path)
}

func (p *PrintDirOption) NPath() int {
	if p.Paths == nil {
		return -1
	}
	return len(p.Paths)
}

func NewPrintDirOption() *PrintDirOption {
	return &PrintDirOption{
		Depth:     0,
		OutOpt:    PListView,
		FieldFlag: PFieldDefault,
		// SortOpt:
		// FiltOpt:,
		Ignore: DefaultIgnoreFn,
		Root:   "",
		Paths:  nil,
	}
}

type PDViewFlag int

const (
	// PListView is the option of list view using in PrintDir
	PListView PDViewFlag = 1 << iota // 1 << 0 which is 00000001
	// PListExtendView is the option of list view icluding extend attributes using in PrintDir
	PListExtendView
	// PTreeView is the option of tree view using in PrintDir
	PTreeView
	// PTreeExtendView is the option of tree view icluding extend atrribute using in PrintDir
	PTreeExtendView
	// PLevelView is the option of level view using in PrintDir
	PLevelView
	// PLevelExtendView is the option of level view icluding extend attributes using in PrintDir
	PLevelExtendView
	// PTableView is the option of table view using in PrintDir
	PTableView
	// PTableView is the option of table view icluding extend attributes using in PrintDir
	PTableExtendView
	// PClassifyView display type indicator by file names (like as `exa -F` or `exa --classify`) in PrintDir
	PClassifyView
	// PListTreeView is the option of combining list & tree view using in PrintDir
	PListTreeView = PListView | PTreeView
	// PListTreeExtendView is the option of combining list & tree view including extend attribute using in PrintDir
	PListTreeExtendView = PListView | PTreeExtendView
)

// PDSortOption defines sorting way view of PrintDir
//
// Defaut:
//  increasing sort by lower name of path
type PDSortOption struct {
	IsSort  bool
	Reverse bool
	SortWay PDSortFlag
}

type PDSortFlag int

const (
	PDSort PDSortFlag = 1 << iota
	PDSortReverse
	pdSortKeyINode
	pdSortKeyLinks
	pdSortKeySize
	pdSortKeyBlocks
	pdSortKeyMTime
	pdSortKeyCTime
	pdSortKeyATime
	pdSortKeyName

	PDSortByINode  = PDSort | pdSortKeyINode
	PDSortByINodeR = PDSortByINode | PDSortReverse

	PDSortByLinks  = PDSort | pdSortKeyLinks
	PDSortByLinksR = PDSortByLinks | PDSortReverse

	PDSortBySize  = PDSort | pdSortKeySize
	PDSortBySizeR = PDSortBySize | PDSortReverse

	PDSortByBlocks  = PDSort | pdSortKeyBlocks
	PDSortByBlocksR = PDSortByBlocks | PDSortReverse

	PDSortByMTime  = PDSort | pdSortKeyMTime
	PDSortByMTimeR = PDSortByMTime | PDSortReverse

	PDSortByCTime  = PDSort | pdSortKeyCTime
	PDSortByCTimeR = PDSortByCTime | PDSortReverse

	PDSortByATime  = PDSort | pdSortKeyATime
	PDSortByATimeR = PDSortByATime | PDSortReverse

	PDSortByName  = PDSort | pdSortKeyName
	PDSortByNameR = PDSortByName | PDSortReverse
)

var (
	// sortedFields = []string{"size", "modified", "accessed", "created", "name"}

	sortByField = map[PDSortFlag]FilesBy{
		PDSortByINode:   byINode,
		PDSortByINodeR:  byINodeR,
		PDSortByLinks:   byLinks,
		PDSortByLinksR:  byLinksR,
		PDSortBySize:    bySize,
		PDSortBySizeR:   bySizeR,
		PDSortByBlocks:  byBlocks,
		PDSortByBlocksR: byBlocksR,
		PDSortByMTime:   byMTime,
		PDSortByMTimeR:  byMTimeR,
		PDSortByATime:   byATime,
		PDSortByATimeR:  byATimeR,
		PDSortByCTime:   byCTime,
		PDSortByCTimeR:  byCTimeR,
		PDSortByName:    byName,
		PDSortByNameR:   byNameR,
	}
)

type PDFiltFlag int

const (
	PDFiltNoEmptyDir = 1 << iota
	PDFiltJustDirs
	PDFiltJustFiles
	PDFiltJustDirsButNoEmpty     = PDFiltNoEmptyDir | PDFiltJustDirs
	PDFiltJustFilesButNoEmptyDir = PDFiltJustFiles
)

type PDFilterOption struct {
	IsFilt  bool
	FiltWay PDFiltFlag
}
