package filetree

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/shyang107/paw/cast"
)

var (
	// NoColor check from the type of terminal and
	// determine output to terminal in color (`true`) or not (`false`)
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))

	// LSColorsFileKindDesc ...
	LSColorsFileKindDesc = map[string]string{
		"di": "directory",
		"fi": "file",
		"ln": "symbolic link",
		"pi": "fifo file",
		"so": "socket file",
		"bd": "block (buffered) special file",
		"cd": "character (unbuffered) special file",
		"or": "symbolic link pointing to a non-existent file (orphan)",
		"mi": "non-existent file pointed to by a symbolic link (visible when you type ls -l)",
		"ex": "file which is executable (ie. has 'x' set in permissions)",
	}

	// LSColors = make(map[string]string) is LS_COLORS code according to
	// extention of file
	LSColors = map[string][]color.Attribute{
		"bd": []color.Attribute{38, 5, 68},
		"ca": []color.Attribute{38, 5, 17},
		"cd": []color.Attribute{38, 5, 113, 1},
		"di": []color.Attribute{38, 5, 30},
		"do": []color.Attribute{38, 5, 127},
		"ex": []color.Attribute{38, 5, 208, 1},
		"pi": []color.Attribute{38, 5, 126},
		"fi": []color.Attribute{0},
		// "ln":                   []color.Attribute{target},
		"mh":                   []color.Attribute{38, 5, 222, 1},
		"no":                   []color.Attribute{0},
		"or":                   []color.Attribute{48, 5, 196, 38, 5, 232, 1},
		"ow":                   []color.Attribute{38, 5, 220, 1},
		"sg":                   []color.Attribute{48, 5, 3, 38, 5, 0},
		"su":                   []color.Attribute{38, 5, 220, 1, 3, 100, 1},
		"so":                   []color.Attribute{38, 5, 197},
		"st":                   []color.Attribute{38, 5, 86, 48, 5, 234},
		"tw":                   []color.Attribute{48, 5, 235, 38, 5, 139, 3},
		"LS_COLORS":            []color.Attribute{48, 5, 89, 38, 5, 197, 1, 3, 4, 7},
		"README":               []color.Attribute{38, 5, 220, 1},
		"README.rst":           []color.Attribute{38, 5, 220, 1},
		"README.md":            []color.Attribute{38, 5, 220, 1},
		"LICENSE":              []color.Attribute{38, 5, 220, 1},
		"COPYING":              []color.Attribute{38, 5, 220, 1},
		"INSTALL":              []color.Attribute{38, 5, 220, 1},
		"COPYRIGHT":            []color.Attribute{38, 5, 220, 1},
		"AUTHORS":              []color.Attribute{38, 5, 220, 1},
		"HISTORY":              []color.Attribute{38, 5, 220, 1},
		"CONTRIBUTORS":         []color.Attribute{38, 5, 220, 1},
		"PATENTS":              []color.Attribute{38, 5, 220, 1},
		"VERSION":              []color.Attribute{38, 5, 220, 1},
		"NOTICE":               []color.Attribute{38, 5, 220, 1},
		"CHANGES":              []color.Attribute{38, 5, 220, 1},
		".log":                 []color.Attribute{38, 5, 190},
		".txt":                 []color.Attribute{38, 5, 253},
		".etx":                 []color.Attribute{38, 5, 184},
		".info":                []color.Attribute{38, 5, 184},
		".markdown":            []color.Attribute{38, 5, 184},
		".md":                  []color.Attribute{38, 5, 184},
		".mkd":                 []color.Attribute{38, 5, 184},
		".nfo":                 []color.Attribute{38, 5, 184},
		".pod":                 []color.Attribute{38, 5, 184},
		".rst":                 []color.Attribute{38, 5, 184},
		".tex":                 []color.Attribute{38, 5, 184},
		".textile":             []color.Attribute{38, 5, 184},
		".bib":                 []color.Attribute{38, 5, 178},
		".json":                []color.Attribute{38, 5, 178},
		".jsonl":               []color.Attribute{38, 5, 178},
		".jsonnet":             []color.Attribute{38, 5, 178},
		".libsonnet":           []color.Attribute{38, 5, 142},
		".ndjson":              []color.Attribute{38, 5, 178},
		".msg":                 []color.Attribute{38, 5, 178},
		".pgn":                 []color.Attribute{38, 5, 178},
		".rss":                 []color.Attribute{38, 5, 178},
		".xml":                 []color.Attribute{38, 5, 178},
		".fxml":                []color.Attribute{38, 5, 178},
		".toml":                []color.Attribute{38, 5, 178},
		".yaml":                []color.Attribute{38, 5, 178},
		".yml":                 []color.Attribute{38, 5, 178},
		".RData":               []color.Attribute{38, 5, 178},
		".rdata":               []color.Attribute{38, 5, 178},
		".xsd":                 []color.Attribute{38, 5, 178},
		".dtd":                 []color.Attribute{38, 5, 178},
		".sgml":                []color.Attribute{38, 5, 178},
		".rng":                 []color.Attribute{38, 5, 178},
		".rnc":                 []color.Attribute{38, 5, 178},
		".cbr":                 []color.Attribute{38, 5, 141},
		".cbz":                 []color.Attribute{38, 5, 141},
		".chm":                 []color.Attribute{38, 5, 141},
		".djvu":                []color.Attribute{38, 5, 141},
		".pdf":                 []color.Attribute{38, 5, 141},
		".PDF":                 []color.Attribute{38, 5, 141},
		".mobi":                []color.Attribute{38, 5, 141},
		".epub":                []color.Attribute{38, 5, 141},
		".docm":                []color.Attribute{38, 5, 111, 4},
		".doc":                 []color.Attribute{38, 5, 111},
		".docx":                []color.Attribute{38, 5, 111},
		".odb":                 []color.Attribute{38, 5, 111},
		".odt":                 []color.Attribute{38, 5, 111},
		".rtf":                 []color.Attribute{38, 5, 111},
		".odp":                 []color.Attribute{38, 5, 166},
		".pps":                 []color.Attribute{38, 5, 166},
		".ppt":                 []color.Attribute{38, 5, 166},
		".pptx":                []color.Attribute{38, 5, 166},
		".ppts":                []color.Attribute{38, 5, 166},
		".pptxm":               []color.Attribute{38, 5, 166, 4},
		".pptsm":               []color.Attribute{38, 5, 166, 4},
		".csv":                 []color.Attribute{38, 5, 78},
		".tsv":                 []color.Attribute{38, 5, 78},
		".ods":                 []color.Attribute{38, 5, 112},
		".xla":                 []color.Attribute{38, 5, 76},
		".xls":                 []color.Attribute{38, 5, 112},
		".xlsx":                []color.Attribute{38, 5, 112},
		".xlsxm":               []color.Attribute{38, 5, 112, 4},
		".xltm":                []color.Attribute{38, 5, 73, 4},
		".xltx":                []color.Attribute{38, 5, 73},
		".pages":               []color.Attribute{38, 5, 111},
		".numbers":             []color.Attribute{38, 5, 112},
		".key":                 []color.Attribute{38, 5, 166},
		"config":               []color.Attribute{1},
		"cfg":                  []color.Attribute{1},
		"conf":                 []color.Attribute{1},
		"rc":                   []color.Attribute{1},
		"authorized_keys":      []color.Attribute{1},
		"known_hosts":          []color.Attribute{1},
		".ini":                 []color.Attribute{1},
		".plist":               []color.Attribute{1},
		".viminfo":             []color.Attribute{1},
		".pcf":                 []color.Attribute{1},
		".psf":                 []color.Attribute{1},
		".hidden-color-scheme": []color.Attribute{1},
		".hidden-tmTheme":      []color.Attribute{1},
		".last-run":            []color.Attribute{1},
		".merged-ca-bundle":    []color.Attribute{1},
		".sublime-build":       []color.Attribute{1},
		".sublime-commands":    []color.Attribute{1},
		".sublime-keymap":      []color.Attribute{1},
		".sublime-settings":    []color.Attribute{1},
		".sublime-snippet":     []color.Attribute{1},
		".sublime-project":     []color.Attribute{1},
		".sublime-workspace":   []color.Attribute{1},
		".tmTheme":             []color.Attribute{1},
		".user-ca-bundle":      []color.Attribute{1},
		".epf":                 []color.Attribute{1},
		".git":                 []color.Attribute{38, 5, 197},
		".gitignore":           []color.Attribute{38, 5, 240},
		".gitattributes":       []color.Attribute{38, 5, 240},
		".gitmodules":          []color.Attribute{38, 5, 240},
		".awk":                 []color.Attribute{38, 5, 172},
		".bash":                []color.Attribute{38, 5, 172},
		".bat":                 []color.Attribute{38, 5, 172},
		".BAT":                 []color.Attribute{38, 5, 172},
		".sed":                 []color.Attribute{38, 5, 172},
		".sh":                  []color.Attribute{38, 5, 172},
		".zsh":                 []color.Attribute{38, 5, 172},
		".vim":                 []color.Attribute{38, 5, 172},
		".kak":                 []color.Attribute{38, 5, 172},
		".ahk":                 []color.Attribute{38, 5, 41},
		".py":                  []color.Attribute{38, 5, 41},
		".ipynb":               []color.Attribute{38, 5, 41},
		".rb":                  []color.Attribute{38, 5, 41},
		".gemspec":             []color.Attribute{38, 5, 41},
		".pl":                  []color.Attribute{38, 5, 208},
		".PL":                  []color.Attribute{38, 5, 160},
		".t":                   []color.Attribute{38, 5, 114},
		".msql":                []color.Attribute{38, 5, 222},
		".mysql":               []color.Attribute{38, 5, 222},
		".pgsql":               []color.Attribute{38, 5, 222},
		".sql":                 []color.Attribute{38, 5, 222},
		".tcl":                 []color.Attribute{38, 5, 64, 1},
		".r":                   []color.Attribute{38, 5, 49},
		".R":                   []color.Attribute{38, 5, 49},
		".gs":                  []color.Attribute{38, 5, 81},
		".clj":                 []color.Attribute{38, 5, 41},
		".cljs":                []color.Attribute{38, 5, 41},
		".cljc":                []color.Attribute{38, 5, 41},
		".cljw":                []color.Attribute{38, 5, 41},
		".scala":               []color.Attribute{38, 5, 41},
		".sc":                  []color.Attribute{38, 5, 41},
		".dart":                []color.Attribute{38, 5, 51},
		".asm":                 []color.Attribute{38, 5, 81},
		".cl":                  []color.Attribute{38, 5, 81},
		".lisp":                []color.Attribute{38, 5, 81},
		".rkt":                 []color.Attribute{38, 5, 81},
		".lua":                 []color.Attribute{38, 5, 81},
		".moon":                []color.Attribute{38, 5, 81},
		".c":                   []color.Attribute{38, 5, 81},
		".C":                   []color.Attribute{38, 5, 81},
		".h":                   []color.Attribute{38, 5, 110},
		".H":                   []color.Attribute{38, 5, 110},
		".tcc":                 []color.Attribute{38, 5, 110},
		".c++":                 []color.Attribute{38, 5, 81},
		".h++":                 []color.Attribute{38, 5, 110},
		".hpp":                 []color.Attribute{38, 5, 110},
		".hxx":                 []color.Attribute{38, 5, 110},
		".ii":                  []color.Attribute{38, 5, 110},
		".M":                   []color.Attribute{38, 5, 110},
		".m":                   []color.Attribute{38, 5, 110},
		".cc":                  []color.Attribute{38, 5, 81},
		".cs":                  []color.Attribute{38, 5, 81},
		".cp":                  []color.Attribute{38, 5, 81},
		".cpp":                 []color.Attribute{38, 5, 81},
		".cxx":                 []color.Attribute{38, 5, 81},
		".cr":                  []color.Attribute{38, 5, 81},
		".go":                  []color.Attribute{38, 5, 81},
		".f":                   []color.Attribute{38, 5, 81},
		".F":                   []color.Attribute{38, 5, 81},
		".for":                 []color.Attribute{38, 5, 81},
		".ftn":                 []color.Attribute{38, 5, 81},
		".f90":                 []color.Attribute{38, 5, 81},
		".F90":                 []color.Attribute{38, 5, 81},
		".f95":                 []color.Attribute{38, 5, 81},
		".F95":                 []color.Attribute{38, 5, 81},
		".f03":                 []color.Attribute{38, 5, 81},
		".F03":                 []color.Attribute{38, 5, 81},
		".f08":                 []color.Attribute{38, 5, 81},
		".F08":                 []color.Attribute{38, 5, 81},
		".nim":                 []color.Attribute{38, 5, 81},
		".nimble":              []color.Attribute{38, 5, 81},
		".s":                   []color.Attribute{38, 5, 110},
		".S":                   []color.Attribute{38, 5, 110},
		".rs":                  []color.Attribute{38, 5, 81},
		".scpt":                []color.Attribute{38, 5, 219},
		".swift":               []color.Attribute{38, 5, 219},
		".sx":                  []color.Attribute{38, 5, 81},
		".vala":                []color.Attribute{38, 5, 81},
		".vapi":                []color.Attribute{38, 5, 81},
		".hi":                  []color.Attribute{38, 5, 110},
		".hs":                  []color.Attribute{38, 5, 81},
		".lhs":                 []color.Attribute{38, 5, 81},
		".agda":                []color.Attribute{38, 5, 81},
		".lagda":               []color.Attribute{38, 5, 81},
		".lagda.tex":           []color.Attribute{38, 5, 81},
		".lagda.rst":           []color.Attribute{38, 5, 81},
		".lagda.md":            []color.Attribute{38, 5, 81},
		".agdai":               []color.Attribute{38, 5, 110},
		".zig":                 []color.Attribute{38, 5, 81},
		".v":                   []color.Attribute{38, 5, 81},
		".pyc":                 []color.Attribute{38, 5, 240},
		".tf":                  []color.Attribute{38, 5, 168},
		".tfstate":             []color.Attribute{38, 5, 168},
		".tfvars":              []color.Attribute{38, 5, 168},
		".css":                 []color.Attribute{38, 5, 125, 1},
		".less":                []color.Attribute{38, 5, 125, 1},
		".sass":                []color.Attribute{38, 5, 125, 1},
		".scss":                []color.Attribute{38, 5, 125, 1},
		".htm":                 []color.Attribute{38, 5, 125, 1},
		".html":                []color.Attribute{38, 5, 125, 1},
		".jhtm":                []color.Attribute{38, 5, 125, 1},
		".mht":                 []color.Attribute{38, 5, 125, 1},
		".eml":                 []color.Attribute{38, 5, 125, 1},
		".mustache":            []color.Attribute{38, 5, 125, 1},
		".coffee":              []color.Attribute{38, 5, 074, 1},
		".java":                []color.Attribute{38, 5, 074, 1},
		".js":                  []color.Attribute{38, 5, 074, 1},
		".mjs":                 []color.Attribute{38, 5, 074, 1},
		".jsm":                 []color.Attribute{38, 5, 074, 1},
		".jsp":                 []color.Attribute{38, 5, 074, 1},
		".php":                 []color.Attribute{38, 5, 81},
		".ctp":                 []color.Attribute{38, 5, 81},
		".twig":                []color.Attribute{38, 5, 81},
		".vb":                  []color.Attribute{38, 5, 81},
		".vba":                 []color.Attribute{38, 5, 81},
		".vbs":                 []color.Attribute{38, 5, 81},
		"Dockerfile":           []color.Attribute{38, 5, 155},
		".dockerignore":        []color.Attribute{38, 5, 240},
		"Makefile":             []color.Attribute{38, 5, 155},
		"MANIFEST":             []color.Attribute{38, 5, 243},
		"pm_to_blib":           []color.Attribute{38, 5, 240},
		".nix":                 []color.Attribute{38, 5, 155},
		".dhall":               []color.Attribute{38, 5, 178},
		".rake":                []color.Attribute{38, 5, 155},
		".am":                  []color.Attribute{38, 5, 242},
		".in":                  []color.Attribute{38, 5, 242},
		".hin":                 []color.Attribute{38, 5, 242},
		".scan":                []color.Attribute{38, 5, 242},
		".m4":                  []color.Attribute{38, 5, 242},
		".old":                 []color.Attribute{38, 5, 242},
		".out":                 []color.Attribute{38, 5, 242},
		".SKIP":                []color.Attribute{38, 5, 244},
		".diff":                []color.Attribute{48, 5, 197, 38, 5, 232},
		".patch":               []color.Attribute{48, 5, 197, 38, 5, 232, 1},
		".bmp":                 []color.Attribute{38, 5, 97},
		".dicom":               []color.Attribute{38, 5, 97},
		".tiff":                []color.Attribute{38, 5, 97},
		".tif":                 []color.Attribute{38, 5, 97},
		".TIFF":                []color.Attribute{38, 5, 97},
		".cdr":                 []color.Attribute{38, 5, 97},
		".flif":                []color.Attribute{38, 5, 97},
		".gif":                 []color.Attribute{38, 5, 97},
		".icns":                []color.Attribute{38, 5, 97},
		".ico":                 []color.Attribute{38, 5, 97},
		".jpeg":                []color.Attribute{38, 5, 97},
		".JPG":                 []color.Attribute{38, 5, 97},
		".jpg":                 []color.Attribute{38, 5, 97},
		".nth":                 []color.Attribute{38, 5, 97},
		".png":                 []color.Attribute{38, 5, 97},
		".psd":                 []color.Attribute{38, 5, 97},
		".pxd":                 []color.Attribute{38, 5, 97},
		".pxm":                 []color.Attribute{38, 5, 97},
		".xpm":                 []color.Attribute{38, 5, 97},
		".webp":                []color.Attribute{38, 5, 97},
		".ai":                  []color.Attribute{38, 5, 99},
		".eps":                 []color.Attribute{38, 5, 99},
		".epsf":                []color.Attribute{38, 5, 99},
		".drw":                 []color.Attribute{38, 5, 99},
		".ps":                  []color.Attribute{38, 5, 99},
		".svg":                 []color.Attribute{38, 5, 99},
		".avi":                 []color.Attribute{38, 5, 114},
		".divx":                []color.Attribute{38, 5, 114},
		".IFO":                 []color.Attribute{38, 5, 114},
		".m2v":                 []color.Attribute{38, 5, 114},
		".m4v":                 []color.Attribute{38, 5, 114},
		".mkv":                 []color.Attribute{38, 5, 114},
		".MOV":                 []color.Attribute{38, 5, 114},
		".mov":                 []color.Attribute{38, 5, 114},
		".mp4":                 []color.Attribute{38, 5, 114},
		".mpeg":                []color.Attribute{38, 5, 114},
		".mpg":                 []color.Attribute{38, 5, 114},
		".ogm":                 []color.Attribute{38, 5, 114},
		".rmvb":                []color.Attribute{38, 5, 114},
		".sample":              []color.Attribute{38, 5, 114},
		".wmv":                 []color.Attribute{38, 5, 114},
		".3g2":                 []color.Attribute{38, 5, 115},
		".3gp":                 []color.Attribute{38, 5, 115},
		".gp3":                 []color.Attribute{38, 5, 115},
		".webm":                []color.Attribute{38, 5, 115},
		".gp4":                 []color.Attribute{38, 5, 115},
		".asf":                 []color.Attribute{38, 5, 115},
		".flv":                 []color.Attribute{38, 5, 115},
		".ts":                  []color.Attribute{38, 5, 115},
		".ogv":                 []color.Attribute{38, 5, 115},
		".f4v":                 []color.Attribute{38, 5, 115},
		".VOB":                 []color.Attribute{38, 5, 115, 1},
		".vob":                 []color.Attribute{38, 5, 115, 1},
		".ass":                 []color.Attribute{38, 5, 117},
		".srt":                 []color.Attribute{38, 5, 117},
		".ssa":                 []color.Attribute{38, 5, 117},
		".sub":                 []color.Attribute{38, 5, 117},
		".sup":                 []color.Attribute{38, 5, 117},
		".vtt":                 []color.Attribute{38, 5, 117},
		".3ga":                 []color.Attribute{38, 5, 137, 1},
		".S3M":                 []color.Attribute{38, 5, 137, 1},
		".aac":                 []color.Attribute{38, 5, 137, 1},
		".amr":                 []color.Attribute{38, 5, 137, 1},
		".au":                  []color.Attribute{38, 5, 137, 1},
		".caf":                 []color.Attribute{38, 5, 137, 1},
		".dat":                 []color.Attribute{38, 5, 137, 1},
		".dts":                 []color.Attribute{38, 5, 137, 1},
		".fcm":                 []color.Attribute{38, 5, 137, 1},
		".m4a":                 []color.Attribute{38, 5, 137, 1},
		".mod":                 []color.Attribute{38, 5, 137, 1},
		".mp3":                 []color.Attribute{38, 5, 137, 1},
		".mp4a":                []color.Attribute{38, 5, 137, 1},
		".oga":                 []color.Attribute{38, 5, 137, 1},
		".ogg":                 []color.Attribute{38, 5, 137, 1},
		".opus":                []color.Attribute{38, 5, 137, 1},
		".s3m":                 []color.Attribute{38, 5, 137, 1},
		".sid":                 []color.Attribute{38, 5, 137, 1},
		".wma":                 []color.Attribute{38, 5, 137, 1},
		".ape":                 []color.Attribute{38, 5, 136, 1},
		".aiff":                []color.Attribute{38, 5, 136, 1},
		".cda":                 []color.Attribute{38, 5, 136, 1},
		".flac":                []color.Attribute{38, 5, 136, 1},
		".alac":                []color.Attribute{38, 5, 136, 1},
		".mid":                 []color.Attribute{38, 5, 136, 1},
		".midi":                []color.Attribute{38, 5, 136, 1},
		".pcm":                 []color.Attribute{38, 5, 136, 1},
		".wav":                 []color.Attribute{38, 5, 136, 1},
		".wv":                  []color.Attribute{38, 5, 136, 1},
		".wvc":                 []color.Attribute{38, 5, 136, 1},
		".afm":                 []color.Attribute{38, 5, 66},
		".fon":                 []color.Attribute{38, 5, 66},
		".fnt":                 []color.Attribute{38, 5, 66},
		".pfb":                 []color.Attribute{38, 5, 66},
		".pfm":                 []color.Attribute{38, 5, 66},
		".ttf":                 []color.Attribute{38, 5, 66},
		".otf":                 []color.Attribute{38, 5, 66},
		".woff":                []color.Attribute{38, 5, 66},
		".woff2":               []color.Attribute{38, 5, 66},
		".PFA":                 []color.Attribute{38, 5, 66},
		".pfa":                 []color.Attribute{38, 5, 66},
		".7z":                  []color.Attribute{38, 5, 40},
		".a":                   []color.Attribute{38, 5, 40},
		".arj":                 []color.Attribute{38, 5, 40},
		".bz2":                 []color.Attribute{38, 5, 40},
		".cpio":                []color.Attribute{38, 5, 40},
		".gz":                  []color.Attribute{38, 5, 40},
		".lrz":                 []color.Attribute{38, 5, 40},
		".lz":                  []color.Attribute{38, 5, 40},
		".lzma":                []color.Attribute{38, 5, 40},
		".lzo":                 []color.Attribute{38, 5, 40},
		".rar":                 []color.Attribute{38, 5, 40},
		".s7z":                 []color.Attribute{38, 5, 40},
		".sz":                  []color.Attribute{38, 5, 40},
		".tar":                 []color.Attribute{38, 5, 40},
		".tgz":                 []color.Attribute{38, 5, 40},
		".warc":                []color.Attribute{38, 5, 40},
		".WARC":                []color.Attribute{38, 5, 40},
		".xz":                  []color.Attribute{38, 5, 40},
		".z":                   []color.Attribute{38, 5, 40},
		".zip":                 []color.Attribute{38, 5, 40},
		".zipx":                []color.Attribute{38, 5, 40},
		".zoo":                 []color.Attribute{38, 5, 40},
		".zpaq":                []color.Attribute{38, 5, 40},
		".zst":                 []color.Attribute{38, 5, 40},
		".zstd":                []color.Attribute{38, 5, 40},
		".zz":                  []color.Attribute{38, 5, 40},
		".apk":                 []color.Attribute{38, 5, 215},
		".ipa":                 []color.Attribute{38, 5, 215},
		".deb":                 []color.Attribute{38, 5, 215},
		".rpm":                 []color.Attribute{38, 5, 215},
		".jad":                 []color.Attribute{38, 5, 215},
		".jar":                 []color.Attribute{38, 5, 215},
		".cab":                 []color.Attribute{38, 5, 215},
		".pak":                 []color.Attribute{38, 5, 215},
		".pk3":                 []color.Attribute{38, 5, 215},
		".vdf":                 []color.Attribute{38, 5, 215},
		".vpk":                 []color.Attribute{38, 5, 215},
		".bsp":                 []color.Attribute{38, 5, 215},
		".dmg":                 []color.Attribute{38, 5, 215},
		".r[0-9]{0,2}":         []color.Attribute{38, 5, 239},
		".zx[0-9]{0,2}":        []color.Attribute{38, 5, 239},
		".z[0-9]{0,2}":         []color.Attribute{38, 5, 239},
		".part":                []color.Attribute{38, 5, 239},
		".iso":                 []color.Attribute{38, 5, 124},
		".bin":                 []color.Attribute{38, 5, 124},
		".nrg":                 []color.Attribute{38, 5, 124},
		".qcow":                []color.Attribute{38, 5, 124},
		".sparseimage":         []color.Attribute{38, 5, 124},
		".toast":               []color.Attribute{38, 5, 124},
		".vcd":                 []color.Attribute{38, 5, 124},
		".vmdk":                []color.Attribute{38, 5, 124},
		".accdb":               []color.Attribute{38, 5, 60},
		".accde":               []color.Attribute{38, 5, 60},
		".accdr":               []color.Attribute{38, 5, 60},
		".accdt":               []color.Attribute{38, 5, 60},
		".db":                  []color.Attribute{38, 5, 60},
		".fmp12":               []color.Attribute{38, 5, 60},
		".fp7":                 []color.Attribute{38, 5, 60},
		".localstorage":        []color.Attribute{38, 5, 60},
		".mdb":                 []color.Attribute{38, 5, 60},
		".mde":                 []color.Attribute{38, 5, 60},
		".sqlite":              []color.Attribute{38, 5, 60},
		".typelib":             []color.Attribute{38, 5, 60},
		".nc":                  []color.Attribute{38, 5, 60},
		".pacnew":              []color.Attribute{38, 5, 33},
		".un~":                 []color.Attribute{38, 5, 241},
		".orig":                []color.Attribute{38, 5, 241},
		".BUP":                 []color.Attribute{38, 5, 241},
		".bak":                 []color.Attribute{38, 5, 241},
		".o":                   []color.Attribute{38, 5, 241},
		"core":                 []color.Attribute{38, 5, 241},
		".mdump":               []color.Attribute{38, 5, 241},
		".rlib":                []color.Attribute{38, 5, 241},
		".dll":                 []color.Attribute{38, 5, 241},
		".swp":                 []color.Attribute{38, 5, 244},
		".swo":                 []color.Attribute{38, 5, 244},
		".tmp":                 []color.Attribute{38, 5, 244},
		".sassc":               []color.Attribute{38, 5, 244},
		".pid":                 []color.Attribute{38, 5, 248},
		".state":               []color.Attribute{38, 5, 248},
		"lockfile":             []color.Attribute{38, 5, 248},
		"lock":                 []color.Attribute{38, 5, 248},
		".err":                 []color.Attribute{38, 5, 160, 1},
		".error":               []color.Attribute{38, 5, 160, 1},
		".stderr":              []color.Attribute{38, 5, 160, 1},
		".aria2":               []color.Attribute{38, 5, 241},
		".dump":                []color.Attribute{38, 5, 241},
		".stackdump":           []color.Attribute{38, 5, 241},
		".zcompdump":           []color.Attribute{38, 5, 241},
		".zwc":                 []color.Attribute{38, 5, 241},
		".pcap":                []color.Attribute{38, 5, 29},
		".cap":                 []color.Attribute{38, 5, 29},
		".dmp":                 []color.Attribute{38, 5, 29},
		".DS_Store":            []color.Attribute{38, 5, 239},
		".localized":           []color.Attribute{38, 5, 239},
		".CFUserTextEncoding":  []color.Attribute{38, 5, 239},
		".allow":               []color.Attribute{38, 5, 112},
		".deny":                []color.Attribute{38, 5, 196},
		".service":             []color.Attribute{38, 5, 45},
		"@.service":            []color.Attribute{38, 5, 45},
		".socket":              []color.Attribute{38, 5, 45},
		".swap":                []color.Attribute{38, 5, 45},
		".device":              []color.Attribute{38, 5, 45},
		".mount":               []color.Attribute{38, 5, 45},
		".automount":           []color.Attribute{38, 5, 45},
		".target":              []color.Attribute{38, 5, 45},
		".path":                []color.Attribute{38, 5, 45},
		".timer":               []color.Attribute{38, 5, 45},
		".snapshot":            []color.Attribute{38, 5, 45},
		".application":         []color.Attribute{38, 5, 116},
		".cue":                 []color.Attribute{38, 5, 116},
		".description":         []color.Attribute{38, 5, 116},
		".directory":           []color.Attribute{38, 5, 116},
		".m3u":                 []color.Attribute{38, 5, 116},
		".m3u8":                []color.Attribute{38, 5, 116},
		".md5":                 []color.Attribute{38, 5, 116},
		".properties":          []color.Attribute{38, 5, 116},
		".sfv":                 []color.Attribute{38, 5, 116},
		".theme":               []color.Attribute{38, 5, 116},
		".torrent":             []color.Attribute{38, 5, 116},
		".urlview":             []color.Attribute{38, 5, 116},
		".webloc":              []color.Attribute{38, 5, 116},
		".lnk":                 []color.Attribute{38, 5, 39},
		"CodeResources":        []color.Attribute{38, 5, 239},
		"PkgInfo":              []color.Attribute{38, 5, 239},
		".nib":                 []color.Attribute{38, 5, 57},
		".car":                 []color.Attribute{38, 5, 57},
		".dylib":               []color.Attribute{38, 5, 241},
		".entitlements":        []color.Attribute{1},
		".pbxproj":             []color.Attribute{1},
		".strings":             []color.Attribute{1},
		".storyboard":          []color.Attribute{38, 5, 196},
		".xcconfig":            []color.Attribute{1},
		".xcsettings":          []color.Attribute{1},
		".xcuserstate":         []color.Attribute{1},
		".xcworkspacedata":     []color.Attribute{1},
		".xib":                 []color.Attribute{38, 5, 208},
		".asc":                 []color.Attribute{38, 5, 192, 3},
		".bfe":                 []color.Attribute{38, 5, 192, 3},
		".enc":                 []color.Attribute{38, 5, 192, 3},
		".gpg":                 []color.Attribute{38, 5, 192, 3},
		".signature":           []color.Attribute{38, 5, 192, 3},
		".sig":                 []color.Attribute{38, 5, 192, 3},
		".p12":                 []color.Attribute{38, 5, 192, 3},
		".pem":                 []color.Attribute{38, 5, 192, 3},
		".pgp":                 []color.Attribute{38, 5, 192, 3},
		".p7s":                 []color.Attribute{38, 5, 192, 3},
		"id_dsa":               []color.Attribute{38, 5, 192, 3},
		"id_rsa":               []color.Attribute{38, 5, 192, 3},
		"id_ecdsa":             []color.Attribute{38, 5, 192, 3},
		"id_ed25519":           []color.Attribute{38, 5, 192, 3},
		".32x":                 []color.Attribute{38, 5, 213},
		".cdi":                 []color.Attribute{38, 5, 213},
		".fm2":                 []color.Attribute{38, 5, 213},
		".rom":                 []color.Attribute{38, 5, 213},
		".sav":                 []color.Attribute{38, 5, 213},
		".st":                  []color.Attribute{38, 5, 213},
		".a00":                 []color.Attribute{38, 5, 213},
		".a52":                 []color.Attribute{38, 5, 213},
		".A64":                 []color.Attribute{38, 5, 213},
		".a64":                 []color.Attribute{38, 5, 213},
		".a78":                 []color.Attribute{38, 5, 213},
		".adf":                 []color.Attribute{38, 5, 213},
		".atr":                 []color.Attribute{38, 5, 213},
		".gb":                  []color.Attribute{38, 5, 213},
		".gba":                 []color.Attribute{38, 5, 213},
		".gbc":                 []color.Attribute{38, 5, 213},
		".gel":                 []color.Attribute{38, 5, 213},
		".gg":                  []color.Attribute{38, 5, 213},
		".ggl":                 []color.Attribute{38, 5, 213},
		".ipk":                 []color.Attribute{38, 5, 213},
		".j64":                 []color.Attribute{38, 5, 213},
		".nds":                 []color.Attribute{38, 5, 213},
		".nes":                 []color.Attribute{38, 5, 213},
		".sms":                 []color.Attribute{38, 5, 213},
		".8xp":                 []color.Attribute{38, 5, 121},
		".8eu":                 []color.Attribute{38, 5, 121},
		".82p":                 []color.Attribute{38, 5, 121},
		".83p":                 []color.Attribute{38, 5, 121},
		".8xe":                 []color.Attribute{38, 5, 121},
		".stl":                 []color.Attribute{38, 5, 216},
		".dwg":                 []color.Attribute{38, 5, 216},
		".ply":                 []color.Attribute{38, 5, 216},
		".wrl":                 []color.Attribute{38, 5, 216},
		".pot":                 []color.Attribute{38, 5, 7},
		".pcb":                 []color.Attribute{38, 5, 7},
		".mm":                  []color.Attribute{38, 5, 7},
		".gbr":                 []color.Attribute{38, 5, 7},
		".scm":                 []color.Attribute{38, 5, 7},
		".xcf":                 []color.Attribute{38, 5, 7},
		".spl":                 []color.Attribute{38, 5, 7},
		".Rproj":               []color.Attribute{38, 5, 11},
		".sis":                 []color.Attribute{38, 5, 7},
		".1p":                  []color.Attribute{38, 5, 7},
		".3p":                  []color.Attribute{38, 5, 7},
		".cnc":                 []color.Attribute{38, 5, 7},
		".def":                 []color.Attribute{38, 5, 7},
		".ex":                  []color.Attribute{38, 5, 7},
		".example":             []color.Attribute{38, 5, 7},
		".feature":             []color.Attribute{38, 5, 7},
		".ger":                 []color.Attribute{38, 5, 7},
		".ics":                 []color.Attribute{38, 5, 7},
		".map":                 []color.Attribute{38, 5, 7},
		".mf":                  []color.Attribute{38, 5, 7},
		".mfasl":               []color.Attribute{38, 5, 7},
		".mi":                  []color.Attribute{38, 5, 7},
		".mtx":                 []color.Attribute{38, 5, 7},
		".pc":                  []color.Attribute{38, 5, 7},
		".pi":                  []color.Attribute{38, 5, 7},
		".plt":                 []color.Attribute{38, 5, 7},
		".pm":                  []color.Attribute{38, 5, 7},
		".rdf":                 []color.Attribute{38, 5, 7},
		".ru":                  []color.Attribute{38, 5, 7},
		".sch":                 []color.Attribute{38, 5, 7},
		".sty":                 []color.Attribute{38, 5, 7},
		".sug":                 []color.Attribute{38, 5, 7},
		".tdy":                 []color.Attribute{38, 5, 7},
		".tfm":                 []color.Attribute{38, 5, 7},
		".tfnt":                []color.Attribute{38, 5, 7},
		".tg":                  []color.Attribute{38, 5, 7},
		".vcard":               []color.Attribute{38, 5, 7},
		".vcf":                 []color.Attribute{38, 5, 7},
		".xln":                 []color.Attribute{38, 5, 7},
		".iml":                 []color.Attribute{38, 5, 166},
	}
	// LS_COLORS code for specific pattern of file extentions
	reExtLSColors = map[*regexp.Regexp][]color.Attribute{
		regexp.MustCompile(`r[0-9]{0,2}$`):  []color.Attribute{38, 5, 239},
		regexp.MustCompile(`zx[0-9]{0,2}$`): []color.Attribute{38, 5, 239},
		regexp.MustCompile(`z[0-9]{0,2}$`):  []color.Attribute{38, 5, 239},
	}
)

func init() {

}

// SetNoColor will set `true` to `NoColor`
func SetNoColor() {
	NoColor = true
}

// DefaultNoColor will resume the default value of `NoColor`
func DefaultNoColor() {
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

func getcolors() {
	colorenv := os.Getenv("LS_COLORS")
	args := strings.Split(colorenv, ":")

	// colors := make(map[string]string)
	// ctypes := make(map[string]string)
	// exts := []string{}
	for _, a := range args {
		// fmt.Printf("%v\t", a)
		kv := strings.Split(a, "=")

		// fmt.Printf("%v\n", kv)
		if len(kv) == 2 {
			LSColors[kv[0]] = getColorAttribute(kv[1])
			// exts = append(exts, kv[0])
		}
	}
	// sort.Strings(exts)
}

func getColorAttribute(code string) []color.Attribute {
	att := []color.Attribute{}
	for _, a := range strings.Split(code, ";") {
		att = append(att, color.Attribute(cast.ToInt(a)))
	}
	return att
}

// FileLSColorString will return the color string of `s` according `fullpath` (xxx.yyy)
func FileLSColorString(fullpath, s string) (string, error) {
	file, ext := getColorExt(fullpath)
	if ext == "«link»" {
		// link, err := os.Readlink(fullpath)
		ext = "ln"
		_, err := filepath.EvalSymlinks(fullpath)
		if err != nil {
			ext = "or"
		}
		// else {
		// 	file, ext = getColorExt(link)
		// }
	}
	switch {
	case NoColor:
		return s, nil
	default:
		if _, ok := LSColors[file]; ok {
			return colorstr(LSColors[file], s), nil
		}
		if _, ok := LSColors[ext]; ok {
			return colorstr(LSColors[ext], s), nil
		}
		for re, att := range reExtLSColors {
			if re.MatchString(file) {
				return colorstr(att, s), nil
			}
		}
		return colorstr(LSColors["no"], s), nil
	}
}

// KindLSColorString will colorful string `s` using key `kind`
func KindLSColorString(kind, s string) string {
	att, ok := LSColors[kind]
	if !ok {
		att = LSColors["fi"]
	}
	return colorstr(att, s)
}

func colorstr(att []color.Attribute, s string) string {
	cs := color.New(att...)
	return cs.Sprint(s)
}

// getColorExt will return the color key of extention from `fullpath`
func getColorExt(fullpath string) (file, ext string) {
	fi, err := os.Lstat(fullpath)
	if err != nil {
		return "", "no"
	}
	mode := fi.Mode()
	sperm := fmt.Sprintf("%v", mode)
	switch {
	case mode.IsDir(): // d: is a directory 資料夾模式
		ext = "di" // di = directory
	case mode&os.ModeSymlink != 0: // L: symbolic link 象徵性的關聯
		ext = "«link»"
		// ext = "ln"
		// link, err := filepath.EvalSymlinks(fullpath)
		// if err != nil {
		// 	ext = "or"
		// } else {
		// 	_, ext = getColorExt(link)
		// }
	// } else { // mi = non-existent file pointed to by a symbolic link (visible when you type ls -l)
	// 	ext = "mi"
	// }
	case mode&os.ModeSocket != 0: // S: Unix domain socket Unix 主機 socket
		ext = "so" // so = socket file
	case mode&os.ModeNamedPipe != 0:
		ext = "pi" //pi = fifo file
	// case mode&os.ModeDevice != 0:
	// 	ext = ""
	// bd = block (buffered) special file
	case mode&os.ModeCharDevice != 0:
		// cd = character (unbuffered) special file
		ext = "cd"
	case mode.IsRegular() && !mode.IsDir() && strings.Contains(sperm, "x"):
		// ex = file which is executable (ie. has 'x' set in permissions)
		ext = "ex"
	default: // fi = file
		ext = filepath.Ext(fullpath)
	}

	return filepath.Base(fullpath), ext
}
