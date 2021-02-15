package filetree

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
)

// GStatus represents the current status of a Worktree.
// The key of the map is the path of the file.
type GStatus map[string]*GitFileStatus

// File returns the GitFileStatus for a given path, if the GitFileStatus doesn't
// exists a new GitFileStatus is added to the map using the path as key.
func (s GStatus) File(path string) *GitFileStatus {
	if _, ok := (s)[path]; !ok {
		s[path] = &GitFileStatus{
			Worktree: GitUntracked,
			Staging:  GitUntracked,
		}
	}

	return s[path]
}

// IsUntracked checks if file for given path is 'Untracked'
func (s GStatus) IsUntracked(path string) bool {
	stat, ok := (s)[filepath.ToSlash(path)]
	return ok && stat.Worktree == GitUntracked
}

// IsClean returns true if all the files are in Unmodified status.
func (s GStatus) IsClean() bool {
	for _, status := range s {
		if status.Worktree != GitUnmodified ||
			status.Staging != GitUnmodified {
			return false
		}
	}

	return true
}

func (s GStatus) String() string {
	buf := bytes.NewBuffer(nil)
	for path, status := range s {
		if status.Staging == GitUnmodified && status.Worktree == GitUnmodified {
			continue
		}

		if status.Staging == GitRenamed {
			path = fmt.Sprintf("%s -> %s", path, status.Extra)
		}

		fmt.Fprintf(buf, "%c%c %s\n", status.Staging, status.Worktree, path)
	}

	return buf.String()
}

// GitFileStatus contains the status of a file in the worktree
type GitFileStatus struct {
	// Staging is the status of a file in the staging area
	Staging GitStatusCode
	// Worktree is the status of a file in the worktree
	Worktree GitStatusCode
	// Extra contains extra information, such as the previous name in a rename
	Extra string
}

// GitStatusCode status code of a file in the Worktree
type GitStatusCode byte

const (
	GitNo                 GitStatusCode = 'X'
	GitUnmodified         GitStatusCode = ' '
	GitUntracked          GitStatusCode = '?'
	GitModified           GitStatusCode = 'M'
	GitAdded              GitStatusCode = 'A'
	GitDeleted            GitStatusCode = 'D'
	GitRenamed            GitStatusCode = 'R'
	GitCopied             GitStatusCode = 'C'
	GitUpdatedButUnmerged GitStatusCode = 'U'
	GitIgnored            GitStatusCode = '!'
	GitChanged            GitStatusCode = 'N'
	GitUnChanged          GitStatusCode = '-' // equl to GitUnmodified
)

func (s GitStatusCode) String() string {
	return string(s)
}

func (s GitStatusCode) Color() *color.Color {
	switch s {
	case GitUnmodified, GitUntracked, GitModified, GitAdded, GitDeleted, GitRenamed, GitCopied, GitUpdatedButUnmerged, GitIgnored, GitChanged, GitUnChanged:
		return cgitmap[s]
	default:
		return cdashp
	}
}

var cgitmap = map[GitStatusCode]*color.Color{
	GitModified:           paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //modified
	GitAdded:              paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //added
	GitDeleted:            paw.NewEXAColor("gd"), //color.New(EXAColors["gd"]...), //deleted
	GitRenamed:            paw.NewEXAColor("gv"), //color.New(EXAColors["gv"]...), //renamed
	GitCopied:             paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //copied
	GitUpdatedButUnmerged: paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //updated but unmerged
	GitUntracked:          paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //untracked
	GitIgnored:            cdashp,                //color.New(EXAColors["-"]...),  //ignored
	GitUnmodified:         cdashp,                //color.New(EXAColors["-"]...),  //ignored
	GitChanged:            paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //untracked
	GitUnChanged:          cdashp,                //color.New(EXAColors["-"]...),  //ignored
	// 'M': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //modified
	// 'A': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //added
	// 'D': paw.NewEXAColor("gd"), //color.New(EXAColors["gd"]...), //deleted
	// 'R': paw.NewEXAColor("gv"), //color.New(EXAColors["gv"]...), //renamed
	// 'C': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //copied
	// 'U': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //updated but unmerged
	// '?': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //untracked
	// 'N': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //untracked
	// '!': cdashp,                //color.New(EXAColors["-"]...),  //ignored
	// '-': cdashp,                //color.New(EXAColors["-"]...),  //ignored
	// ' ': cdashp,                //color.New(EXAColors["-"]...),  //ignored
}

// GitStatus stores git status of `Branch`
// 	NoGit are true : no git
// 	Branch are branch of git
// 	FilesStatus are map[{{ path }}]{{ XY }}
// 		XY are ??, 2 characters string, see also "https://git-scm.com/docs/git-status"
type GitStatus struct {
	NoGit   bool
	repPath string
	// Repository represents a git repository
	// type Repository struct {
	// 	Storer storage.Storer
	// 	// contains filtered or unexported fields
	// }
	// repository *git.Repository
	// Worktree represents a git worktree.
	// type Worktree struct {
	// 	// Filesystem underlying filesystem.
	// 	Filesystem billy.Filesystem
	// 	// External excludes not found in the repository .gitignore
	// 	Excludes []gitignore.Pattern
	// 	// contains filtered or unexported fields
	// }
	// worktree *git.Worktree
	// Status represents the current status of a Worktree. The key of the map is the path of the file.
	// type Status map[string]*GitFileStatus
	branch string
	status GStatus
	// status git.Status
}

func NewGitStatus(repPath string) *GitStatus {
	gs, err := getShortGitStatus(repPath)
	if err != nil {
		return &GitStatus{
			NoGit: true,
		}
	}
	return gs
}

func (g *GitStatus) GetRepositoryPath() string {
	if g.NoGit {
		return ""
	}
	return g.repPath
}

func (g *GitStatus) GetBranch() string {
	if g.NoGit {
		return ""
	}
	return g.branch
}

func (g *GitStatus) GetStatus() GStatus {
	if g.NoGit {
		return nil
	}
	return g.status
}

func (g *GitStatus) SetStatus(gs GStatus) {
	g.status = gs
}

func xy(xy GitStatusCode) GitStatusCode {
	switch xy {
	case GitUntracked:
		return GitChanged
	case GitUnmodified:
		return GitUnChanged
	default:
		return xy
	}
}

func (g *GitStatus) XStaging(relpath string) GitStatusCode {
	if g.NoGit {
		return GitNo
	}
	var x GitStatusCode
	if s, ok := g.status[relpath]; !ok {
		x = GitUnChanged
	} else {
		x = xy(s.Staging)
	}
	return x
}

func (g *GitStatus) XStagingS(relpath string) string {
	return g.XStaging(relpath).String()
}

func (g *GitStatus) XStagingC(relpath string) string {
	x := g.XStaging(relpath)
	if x == GitNo {
		return ""
	}
	return x.Color().Sprint(x.String())
}

func (g *GitStatus) YWorktree(relpath string) GitStatusCode {
	if g.NoGit {
		return GitStatusCode('0')
	}
	var y GitStatusCode
	if s, ok := g.status[relpath]; !ok {
		y = GitUnChanged
	} else {
		y = xy(s.Worktree)
	}
	return y
}

func (g *GitStatus) YWorktreeS(relpath string) string {
	return g.YWorktree(relpath).String()
}

func (g *GitStatus) YWorktreeC(relpath string) string {
	y := g.YWorktree(relpath)
	if y == GitNo {
		return ""
	}
	return y.Color().Sprint(y.String())
}

func (g *GitStatus) XYStatus(relpath string) string {
	return g.XStagingS(relpath) + g.YWorktreeS(relpath)
}

func (g *GitStatus) XYStatusC(relpath string) string {
	return g.XStagingC(relpath) + g.YWorktreeC(relpath)
}

//getShortGitStatus read the git status of the repository located at path
// 	if err != nil : no git
func getShortGitStatus(repPath string) (*GitStatus, error) {
	out, err := execOutput(fmt.Sprintf("git -C %s status -s -b --porcelain --ignored", repPath))
	if err != nil {
		// paw.Logger.Errorf("unable to read git repository status : %s", err.Error())
		return &GitStatus{NoGit: true}, err
	}
	// paw.Logger.WithField("out", out).Debug("git")
	status := parseShort(repPath, out)

	return status, err
}

//It is useful to declare a var instead of a function for testing purpose
var execOutput = func(c string) (io.Reader, error) {
	out, err := exec.Command("/bin/sh", "-c", c).Output()
	return bytes.NewReader(out), err
}

//Parse parses a git status output command
//It is compatible with the short version of the git status command
func parseShort(reppath string, r io.Reader) *GitStatus {
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

	gs := make(GStatus)
	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}
		st := s.Text()
		rfile := st[3:]
		_, file := filepath.Split(rfile)
		gs[rfile] = &GitFileStatus{
			Staging:  GitStatusCode(st[0]),
			Worktree: GitStatusCode(st[1]),
			Extra:    file,
		}
	}
	return &GitStatus{
		NoGit:   false,
		branch:  branch,
		repPath: reppath,
		status:  gs,
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

// const (
// 	Unmodified         StatusCode = ' '
// 	Untracked          StatusCode = '?'
// 	Modified           StatusCode = 'M'
// 	Added              StatusCode = 'A'
// 	Deleted            StatusCode = 'D'
// 	Renamed            StatusCode = 'R'
// 	Copied             StatusCode = 'C'
// 	UpdatedButUnmerged StatusCode = 'U'
// )

// // GitStatus stores git status of `Branch`
// // 	NoGit are true : no git
// // 	Branch are branch of git
// // 	FilesStatus are map[{{ path }}]{{ XY }}
// // 		XY are ??, 2 characters string, see also "https://git-scm.com/docs/git-status"
// type GitStatus struct {
// 	NoGit       bool
// 	Branch      string
// 	FilesStatus map[string]XY // == map[string]XY
// 	// XY are ??
// 	// ' ' = unmodified
// 	// M = modified
// 	// A = added
// 	// D = deleted
// 	// R = renamed
// 	// C = copied
// 	// U = updated but unmerged
// 	//
// 	// Ignored files are not listed, unless --ignored option is in effect, in which case XY are !!.
// 	// X          Y     Meaning
// 	// -------------------------------------------------
// 	//          [AMD]   not updated
// 	// M        [ MD]   updated in index
// 	// A        [ MD]   added to index
// 	// D                deleted from index
// 	// R        [ MD]   renamed in index
// 	// C        [ MD]   copied in index
// 	// [MARC]           index and work tree matches
// 	// [ MARC]     M    work tree changed since index
// 	// [ MARC]     D    deleted in work tree
// 	// [ D]        R    renamed in work tree
// 	// [ D]        C    copied in work tree
// 	// -------------------------------------------------
// 	// D           D    unmerged, both deleted
// 	// A           U    unmerged, added by us
// 	// U           D    unmerged, deleted by them
// 	// U           A    unmerged, added by them
// 	// D           U    unmerged, deleted by us
// 	// A           A    unmerged, both added
// 	// U           U    unmerged, both modified
// 	// -------------------------------------------------
// 	// ?           ?    untracked
// 	// !           !    ignored
// 	// -------------------------------------------------
// }

// type XY []rune

// func (s XY) String() string {
// 	var str string
// 	for _, c := range s {
// 		str += string(c)
// 	}
// 	return str
// }
// func (s XY) Split() (x, y rune) {
// 	x = s[0]
// 	y = s[1]
// 	return x, y
// }

// func ToXY(st string) XY {
// 	return XY{rune(st[0]), rune(st[1])}
// }

// //GetShortGitStatus read the git status of the repository located at path
// // 	if err != nil : no git
// func GetShortGitStatus(path string) (GitStatus, error) {

// 	out, err := execOutput(fmt.Sprintf("git -C %s status -s -b --porcelain", path))
// 	if err != nil {
// 		// paw.Logger.Errorf("unable to read git repository status : %s", err.Error())
// 		return GitStatus{NoGit: true}, err
// 	}

// 	status := ParseShort(path, out)

// 	return status, err
// }

// //It is useful to declare a var instead of a function for testing purpose
// var execOutput = func(c string) (io.Reader, error) {
// 	out, err := exec.Command("/bin/sh", "-c", c).Output()

// 	return bytes.NewReader(out), err
// }

// //Parse parses a git status output command
// //It is compatible with the short version of the git status command
// func ParseShort(root string, r io.Reader) GitStatus {

// 	s := bufio.NewScanner(r)

// 	var branch string
// 	//Extract branch name
// 	for s.Scan() {
// 		//Skip any empty line
// 		if len(s.Text()) < 1 {
// 			continue
// 		}

// 		branch = parseBranch(s.Text())
// 		break
// 	}

// 	fs := make(map[string]XY)
// 	for s.Scan() {
// 		if len(s.Text()) < 1 {
// 			continue
// 		}
// 		st := s.Text()
// 		file := filepath.Join(root, st[3:])
// 		// gstat := strings.Replace(st[:2], " ", "-", -1)
// 		fs[file] = []rune{rune(st[0]), rune(st[1])}
// 	}
// 	return GitStatus{
// 		NoGit:       false,
// 		Branch:      branch,
// 		FilesStatus: fs,
// 	}
// }

// func parseBranch(input string) string {
// 	s := bufio.NewScanner(strings.NewReader(input))
// 	s.Split(bufio.ScanWords)

// 	//check if input is a status branch line output
// 	s.Scan()
// 	if s.Text() != "##" {
// 		return ""
// 	}

// 	//read next word and return the branch name
// 	s.Scan()
// 	b := strings.Split(s.Text(), "...")
// 	return b[0]
// }
