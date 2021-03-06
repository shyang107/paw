package vfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
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
		if (status.Staging == GitUnmodified &&
			status.Worktree == GitUnmodified) ||
			(status.Staging == GitUnChanged &&
				status.Worktree == GitUnChanged) {
			return false
		}
	}

	return true
}

func (s GStatus) String() string {
	buf := new(strings.Builder)
	for path, status := range s {
		if (status.Staging == GitUnmodified &&
			status.Worktree == GitUnmodified) ||
			(status.Staging == GitUnChanged &&
				status.Worktree == GitUnChanged) {
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
	GitNo GitStatusCode = 'X'
	// GitUnmodified in input is replaced with GitUnChanged
	GitUnmodified         GitStatusCode = ' '
	GitUntracked          GitStatusCode = '?'
	GitModified           GitStatusCode = 'M'
	GitAdded              GitStatusCode = 'A'
	GitDeleted            GitStatusCode = 'D'
	GitRenamed            GitStatusCode = 'R'
	GitCopied             GitStatusCode = 'C'
	GitUpdatedButUnmerged GitStatusCode = 'U'
	// gitIgnored in input is replaced with GitIgnored
	gitIgnored   GitStatusCode = '!'
	GitIgnored   GitStatusCode = 'I'
	GitChanged   GitStatusCode = 'N'
	GitUnChanged GitStatusCode = '-' // equl to GitUnmodified
)

func (s GitStatusCode) String() string {
	return string(s)
}

func (s GitStatusCode) Color() *color.Color {
	if c, ok := cgitmap[s]; !ok {
		return cdashp
	} else {
		return c
	}
}

var cgitmap = map[GitStatusCode]*color.Color{
	GitNo:                 cdashp,
	GitUnmodified:         cdashp,
	GitUntracked:          paw.NewEXAColor("gm"),
	GitModified:           paw.NewEXAColor("gm"),
	GitAdded:              paw.NewEXAColor("ga"),
	GitDeleted:            paw.NewEXAColor("gd"),
	GitRenamed:            paw.NewEXAColor("gv"),
	GitCopied:             paw.NewEXAColor("gv"),
	GitUpdatedButUnmerged: paw.NewEXAColor("gt"),
	GitIgnored:            cdashp,
	GitChanged:            paw.NewEXAColor("ga"),
	GitUnChanged:          cdashp,
}

// GitStatus stores git status of `Branch`
// 	NoGit are true : no git
// 	Branch are branch of git
// 	FilesStatus are map[{{ path }}]{{ XY }}
// 		XY are ??, 2 characters string, see also "https://git-scm.com/docs/git-status"
type GitStatus struct {
	NoGit   bool
	repPath string
	head    string
	status  GStatus
	// status git.Status
}

func NewGitStatus(repPath string) *GitStatus {
	gs, err := getShortGitStatus(repPath)
	if err != nil {
		return &GitStatus{
			NoGit: true,
		}
	}
	// paw.Logger.Debug(gs.head)
	if paw.Logger.IsLevelEnabled(logrus.TraceLevel) {
		gs.Dump("NewGitStatus")
	}
	return gs
}

func getSC(sc []GitStatusCode) GitStatusCode {
	if len(sc) == 0 {
		return GitUnmodified
	}
	c := sc[0]
	for i := 1; i < len(sc); i++ {
		if c != sc[i] {
			return GitChanged
		}
		c = sc[i]
	}
	return c
}

func (g *GitStatus) Dump(msg string) {
	if len(msg) > 0 {
		paw.Logger.Debugf("[%v] branch: %v%v", msg, g.head, paw.Caller(1))
	} else {
		paw.Logger.Debug(g.head + paw.Caller(1))
	}

	rps := []string{}
	for rp := range g.status {
		rps = append(rps, rp)
	}

	sort.Slice(rps, func(i, j int) bool {
		return strings.ToLower(rps[i]) < strings.ToLower(rps[j])
	})

	// cdp := color.New(color.FgMagenta).Add(color.Bold)
	// cfp := cfip // color.New(color.FgBlue)
	for i, rp := range rps {
		v := g.status[rp]
		// var crp string
		// if strings.HasSuffix(rp, "/") {
		// 	crp = cdp.Sprintf("%q", rp) + ""
		// } else {
		// 	crp = cfp.Sprintf("%q", rp) + ""
		// }
		paw.Logger.WithFields(logrus.Fields{
			"":   i,
			"rp": rp,
			"X":  v.Staging,
			"Y":  v.Worktree,
			"x":  `"` + v.Extra + `"`,
		}).Trace()
	}
}

func (g *GitStatus) GetRepositoryPath() string {
	if g.NoGit {
		return ""
	}
	return g.repPath
}

func (g *GitStatus) GetHead() string {
	if g.NoGit {
		return ""
	}
	return g.head
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
	if s, ok := g.status[relpath]; !ok {
		return GitUnChanged
	} else {
		return xy(s.Staging)
	}
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
	if s, ok := g.status[relpath]; !ok {
		return GitUnChanged
	} else {
		return xy(s.Worktree)
	}
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

func (g *GitStatus) XY(relpath string) string {
	return g.XStagingS(relpath) + g.YWorktreeS(relpath)
}

func (g *GitStatus) XYc(relpath string) string {
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
	// paw.Logger.WithField("out", out).Trace("git")
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

		// branch = parseBranch(s.Text())
		branch = strings.TrimPrefix(s.Text(), "## ")
		break
	}

	gs := make(GStatus)
	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}
		st := s.Text()
		rfile := st[3:]
		// _, file := filepath.Split(rfile)
		rfs := strings.Split(rfile, PathSeparator)
		nrfs := len(rfs)
		file := rfs[nrfs-1]
		if len(file) == 0 {
			file = rfs[nrfs-2]
			if strings.HasSuffix(rfile, PathSeparator) {
				file += PathSeparator
			}
		}
		x := GitStatusCode(st[0])
		y := GitStatusCode(st[1])
		if x == gitIgnored {
			x = GitIgnored
		}
		if y == gitIgnored {
			y = GitIgnored
		}
		gs[rfile] = &GitFileStatus{
			Staging:  x,
			Worktree: y,
			Extra:    file,
		}
	}
	return &GitStatus{
		NoGit:   false,
		head:    branch,
		repPath: reppath,
		status:  gs,
	}
}

func parseBranch(input string) (branch string) {
	if !strings.HasPrefix(branch, "## ") {
		return ""
	} else {
		return branch[3:]
	}
	// s := bufio.NewScanner(strings.NewReader(input))
	// s.Split(bufio.ScanWords)

	// //check if input is a status branch line output
	// s.Scan()
	// if s.Text() != "##" {
	// 	return ""
	// }

	// //read next word and return the branch name
	// // branch := strings.Split(s.Text(), "...")
	// // return branch[0]
	// for s.Scan() {
	// 	branch += s.Text() + " "
	// }

	// return strings.TrimSpace(branch)
}
