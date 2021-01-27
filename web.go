package paw

import (
	"os/exec"
	"strings"
)

var (
	browserApp = map[string]string{
		// Edge "Microsoft Edge"
		"edge": "Microsoft Edge",
		// Chrome "Google Chrome"
		"chrome": "Google Chrome",
	}
)

var sep = `«»`

// GetTitleAndURL get the `title` and `UTL` of active tab of the current window of `browser`
//    `browser`:
//       "edge" for "Microsoft Edge" (default)
//       "chrome" for "Google Chrome"
func GetTitleAndURL(browser string) (t, u string, err error) {
	browser, ok := browserApp[browser]
	if !ok {
		browser = browserApp["edge"] // default browser
	}
	osa, _ := exec.LookPath("osascript")
	cmd := exec.Command(osa,
		"-e", `tell application "`+browser+`"`,
		"-e", `tell active tab of front window`,
		"-e", `set {t, u} to {title, URL}`,
		"-e", `end tell`,
		"-e", `end tell`,
		"-e", `return t & "`+sep+`" & u`)
	// fmt.Println(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		// log.Warningln(err)
		return "", "", err
	}
	str := strings.TrimSpace(string(out))
	tus := strings.Split(str, sep)
	t, u = strings.TrimSpace(tus[0]), strings.TrimSpace(tus[1])
	return t, u, nil
}

// GetTitle get the `title` of active tab of the current window of `browser`
func GetTitle(browser string) (string, error) {
	t, _, err := GetTitleAndURL(browser)
	return t, err
}

// GetURL get the `URL` of active tab of the current window of `browser`
func GetURL(browser string) (string, error) {
	_, u, err := GetTitleAndURL(browser)
	return u, err
}
