package main

import (
	"github.com/go-git/go-git"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func main() {
	paw.Logger.SetLevel(logrus.TraceLevel)
	reppath := "/Users/shyang/go/src/github.com/shyang107/paw/"
	// gs := filetree.GitStatus{
	// 	NoGit: true,
	// }
	// r, err := git.PlainOpen(reppath)
	r, err := git.PlainOpenWithOptions(reppath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		paw.Logger.Error(err)
	}

	w, err := r.Worktree()
	if err != nil {
		paw.Logger.Error(err)
	}
	ws, err := w.Status()
	if err != nil {
		paw.Logger.Error(err)
	}
	for rpath, xy := range ws {
		paw.Logger.WithFields(logrus.Fields{
			"rpath": rpath,
			"X":     string(xy.Staging),
			"Y":     string(xy.Worktree),
		}).Trace()
	}
	// filename := "filetree/filetree_helper.go"
	filename := "examples/git/"
	var st *git.FileStatus
	if _, ok := ws[filename]; !ok {
		// st = ws.File(filename)
		paw.Logger.WithFields(logrus.Fields{
			"filename": filename, // X
		}).Info()
	}
	st = ws.File(filename)

	paw.Logger.WithFields(logrus.Fields{
		"Staging":  gitStatusString(st.Staging),  // X
		"Worktree": gitStatusString(st.Worktree), //Y
		"Extra":    st.Extra,
	}).Info()
}

func gitStatusString(statusCode git.StatusCode) string {
	return string(statusCode)
	// return fmt.Sprintf("%q", statusCode)
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
