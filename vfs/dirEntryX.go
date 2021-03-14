package vfs

import (
	"io/fs"
	"time"

	"github.com/fatih/color"
)

type FileMode = fs.FileMode
type FileInfo = fs.FileInfo
type DirEntry = fs.DirEntry

// Extended is a interface to get extended attributes from Dir or File
type Extendeder interface {
	Xattibutes() []string
}

type Fielder interface {
	Path() string
	RelPath() string
	RelDir() string
	NameToLink() string
	LinkPath() string
	LSColor() *color.Color
	INode() uint64
	HDLinks() uint64
	Blocks() uint64
	Uid() uint32
	User() string
	Gid() uint32
	Group() string
	Dev() uint64
	DevNumber() (uint32, uint32)
	DevNumberS() string
	AccessedTime() time.Time
	CreatedTime() time.Time
	ModifiedTime() time.Time
	Md5() string
	Git() *GitStatus
	XY() string

	Field(ViewField) string
	FieldC(ViewField) string
	WidthOf(ViewField) int
}

type ISer interface {
	IsLink() bool
	IsFile() bool
	IsCharDev() bool
	IsDev() bool
	IsFIFO() bool
	IsSocket() bool
	IsTemporary() bool
	IsExecOwner() bool
	IsExecGroup() bool
	IsExecOther() bool
	IsExecAny() bool
	IsExecAll() bool
	IsExecutable() bool
}

type DirEntryX interface {
	FileInfo
	DirEntry
	Extendeder
	Fielder
	ISer
}
