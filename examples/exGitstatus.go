package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/shyang107/paw"

	"github.com/shyang107/paw/filetree"
)

func exGitstatus(path string) {
	path, _ = filepath.Abs(path)
	status, err := filetree.GetShortStatus(path)
	if err != nil {
		log.Errorf("unable to read git repository status : %s", err.Error())
	}
	if status.NoGit {
		paw.Logger.Warnln("No git")
	} else {
		fmt.Println("Branch is ", status.Branch)
		for file, status := range status.FilesStatus {
			fmt.Printf("%q|%q\n", file, status)
		}
	}

	fl := filetree.NewFileList(path)
	fl.FindFiles(0, nil)
	fmt.Println(fl.ToListViewString(""))
	dirs := fl.Dirs()
	fm := fl.Map()
	git := fl.GetGitStatus()
	spew.Dump(git)
	for _, dir := range dirs {
		for _, file := range fm[dir] {
			fmt.Println("  ", file, getGit(git, file))
		}
	}
}
func getGit(git filetree.GitStatus, file *filetree.File) string {
	st := "--"
	if xy, ok := git.FilesStatus[file.Path]; ok {
		return xy.String()
	}
	if file.IsDir() {
		gits := getGitSlice(git, file)
		if len(gits) > 0 {
			fmt.Println("> ", file.Path, file.Dir, gits)
			return getGitTag(gits)
		} else {
			return st
		}
	}
	return st

	// return  st
}

func getGitTag(gits []string) string {
	paw.Logger.Info()
	x := getGitTagChar(rune(gits[0][0]))
	y := getGitTagChar(rune(gits[0][1]))
	for i := 1; i < len(gits); i++ {
		c := getGitTagChar(rune(gits[i][0]))
		if c != '-' && x != 'N' {
			x = c
		}
		c = getGitTagChar(rune(gits[i][1]))
		if c != '-' && y != 'N' {
			y = c
		}
	}
	return string(x) + string(y)
}

func getGitTagChar(c rune) rune {
	if c == '?' || c == 'A' {
		return 'N'
	}
	return c
}

func getGitSlice(git filetree.GitStatus, file *filetree.File) []string {
	gits := []string{}
	for k, v := range git.FilesStatus {
		if strings.HasPrefix(k, file.Path) {
			ss := strings.Replace(v.String(), " ", "-", 1)
			gits = append(gits, ss)
		}
	}
	return gits
}
