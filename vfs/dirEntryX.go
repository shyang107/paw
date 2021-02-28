package vfs

import (
	"io/fs"
	"time"

	"github.com/fatih/color"
)

// Extended is a interface to get extended attributes from Dir or File
type Extended interface {
	Xattibutes() []string
}

type DirEntryX interface {
	fs.FileInfo
	fs.DirEntry
	Extended

	Path() string
	RelPath() string
	NameToLink() string
	LinkPath() string
	LSColor() *color.Color
	INode() uint64
	HDLinks() uint64
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
	WidthOf(ViewField) int
	// Field(ViewField, *GitStatus) string
	// FieldC(ViewField, *GitStatus)
}
