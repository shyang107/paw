package filetree

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitStatus stores git status of `Branch`
// 	NoGit are true : no git
// 	Branch are branch of git
// 	FilesStatus are map[{{ path }}]{{ XY }}
// 		XY are ??, 2 characters string, see also "https://git-scm.com/docs/git-status"
type GitStatus struct {
	NoGit       bool
	Branch      string
	FilesStatus map[string]XY // == map[string]XY
	// XY are ??
	// ' ' = unmodified
	// M = modified
	// A = added
	// D = deleted
	// R = renamed
	// C = copied
	// U = updated but unmerged
	//
	// Ignored files are not listed, unless --ignored option is in effect, in which case XY are !!.
	// X          Y     Meaning
	// -------------------------------------------------
	//          [AMD]   not updated
	// M        [ MD]   updated in index
	// A        [ MD]   added to index
	// D                deleted from index
	// R        [ MD]   renamed in index
	// C        [ MD]   copied in index
	// [MARC]           index and work tree matches
	// [ MARC]     M    work tree changed since index
	// [ MARC]     D    deleted in work tree
	// [ D]        R    renamed in work tree
	// [ D]        C    copied in work tree
	// -------------------------------------------------
	// D           D    unmerged, both deleted
	// A           U    unmerged, added by us
	// U           D    unmerged, deleted by them
	// U           A    unmerged, added by them
	// D           U    unmerged, deleted by us
	// A           A    unmerged, both added
	// U           U    unmerged, both modified
	// -------------------------------------------------
	// ?           ?    untracked
	// !           !    ignored
	// -------------------------------------------------
}

type XY []rune

func (s XY) String() string {
	var str string
	for _, c := range s {
		str += string(c)
	}
	return str
}
func (s XY) Split() (x, y rune) {
	x = s[0]
	y = s[1]
	return x, y
}

func ToXY(st string) XY {
	return XY{rune(st[0]), rune(st[1])}
}

//GetShortStatus read the git status of the repository located at path
// 	if err != nil : no git
func GetShortStatus(path string) (GitStatus, error) {
	out, err := execOutput(fmt.Sprintf("git -C %s status -s -b --porcelain", path))
	if err != nil {
		// paw.Logger.Errorf("unable to read git repository status : %s", err.Error())
		return GitStatus{NoGit: true}, err
	}

	status := ParseShort(path, out)

	return status, err
}

//It is useful to declare a var instead of a function for testing purpose
var execOutput = func(c string) (io.Reader, error) {
	out, err := exec.Command("/bin/sh", "-c", c).Output()

	return bytes.NewReader(out), err
}

//Parse parses a git status output command
//It is compatible with the short version of the git status command
func ParseShort(root string, r io.Reader) GitStatus {

	s := bufio.NewScanner(r)

	var branch string
	//Extract branch name
	for s.Scan() {
		//Skip any empty line
		if len(s.Text()) < 1 {
			continue
		}

		branch = parseBranch(s.Text())
		break
	}

	fs := make(map[string]XY)
	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}
		st := s.Text()
		file := filepath.Join(root, st[3:])
		// gstat := strings.Replace(st[:2], " ", "-", -1)
		fs[file] = []rune{rune(st[0]), rune(st[1])}
	}
	return GitStatus{
		NoGit:       false,
		Branch:      branch,
		FilesStatus: fs,
	}
}

func parseBranch(input string) string {
	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(bufio.ScanWords)

	//check if input is a status branch line output
	s.Scan()
	if s.Text() != "##" {
		return ""
	}

	//read next word and return the branch name
	s.Scan()
	b := strings.Split(s.Text(), "...")
	return b[0]
}
