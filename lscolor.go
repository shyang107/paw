package paw

import (
	"os"
	"path/filepath"

	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/shyang107/paw/cast"
	"github.com/sirupsen/logrus"
)

type Color = color.Color
type Attribute = color.Attribute

const (
	Gray0 Attribute = iota + 232
	Gray1
	Gray2
	Gray3
	Gray4
	Gray5
	Gray6
	Gray7
	Gray8
	Gray9
	Gray10
	Gray11
	Gray12
	Gray13
	Gray14
	Gray15
	Gray16
	Gray17
	Gray18
	Gray19
	Gray20
	Gray21
	Gray22
	Gray23
)

var (
	// NoColor check from the type of terminal and
	// determine output to terminal in color (`true`) or not (`false`)
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))

	// LSColorsFileKindDesc ...
	LSColorsFileKindDesc = map[string]string{
		"ca": "file with capability",
		"bd": "block (buffered) device, special file",
		"cd": "character (unbuffered) device, special file",
		"do": "door",
		"di": "directory",
		"ec": "ENDCODE, non-filename text",
		"ex": "file which is executable (ie. has 'x' set in permissions).",
		"fi": "file",
		// "lc":          "LEFTCODE, opening terminal code",
		"ln": "symbolic link",
		// "mi":          "non-existent file pointed to by a symbolic link (visible when you type ls -l)",
		"no":          "normal, global default",
		"or":          "orphan, symbolic link pointing to a non-existent file (orphan)",
		"pi":          "fifo file, named pipe",
		"rc":          "RIGHTCODE, closing terminal code",
		"sg":          "file that is setgid (g+s)",
		"so":          "socket file",
		"st":          "sticky bit set (+t) and not other-writable directory",
		"su":          "file that is setuid (u+s)",
		"tw":          "sticky and other-writable (+t,o+w) directory",
		"ow":          "other-writable (o+w) and not sticky directory",
		"*.extension": "every file using this extension e.g. *.jpg",
	}

	// Grays: 0-23 gradient (color.Attribute) of gray
	// 	 232-255：從黑到白的24階灰度色
	Grays = map[int]Attribute{
		0:  Gray0,
		1:  Gray1,
		2:  Gray2,
		3:  Gray3,
		4:  Gray4,
		5:  Gray5,
		6:  Gray6,
		7:  Gray7,
		8:  Gray8,
		9:  Gray9,
		10: Gray10,
		11: Gray11,
		12: Gray12,
		13: Gray13,
		14: Gray14,
		15: Gray15,
		16: Gray16,
		17: Gray17,
		18: Gray18,
		19: Gray19,
		20: Gray20,
		21: Gray21,
		22: Gray22,
		23: Gray23,
	}
	// GraysI: 0-23 gradient (int) of gray
	// 	232-255：從黑到白的24階灰度色
	GraysI = map[int]int{
		0:  int(Gray0),
		1:  int(Gray1),
		2:  int(Gray2),
		3:  int(Gray3),
		4:  int(Gray4),
		5:  int(Gray5),
		6:  int(Gray6),
		7:  int(Gray7),
		8:  int(Gray8),
		9:  int(Gray9),
		10: int(Gray10),
		11: int(Gray11),
		12: int(Gray12),
		13: int(Gray13),
		14: int(Gray14),
		15: int(Gray15),
		16: int(Gray16),
		17: int(Gray17),
		18: int(Gray18),
		19: int(Gray19),
		20: int(Gray20),
		21: int(Gray21),
		22: int(Gray22),
		23: int(Gray23),
	}

	EXAColorAttributes = map[string][]Attribute{
		"ca": LSColorAttributes["ca"],
		"cd": LSColorAttributes["cd"],
		"di": LSColorAttributes["di"],
		"do": LSColorAttributes["do"],
		"ex": LSColorAttributes["ex"],
		"pi": LSColorAttributes["pi"],
		"fi": LSColorAttributes["fi"],
		"ln": LSColorAttributes["ln"],
		"mh": LSColorAttributes["mh"],
		"no": LSColorAttributes["no"],
		"or": LSColorAttributes["or"],
		"ow": LSColorAttributes["ow"],
		"sg": LSColorAttributes["sg"],
		"su": LSColorAttributes["su"],
		"so": LSColorAttributes["so"],
		"st": LSColorAttributes["st"],
		"bd": LSColorAttributes["bd"],
		"rc": LSColorAttributes["rc"],
		// "ur": LSColorAttributes["ex"],
		"ur": FgColor256A(230).Add(color.Bold),
		//{38, 5, 230, 1}, // user +r bit
		"uw": FgColor256A(209).Add(color.Bold),
		//{38, 5, 209, 1}, // user +w bit
		"ux": FgColor256A(156).Add(color.Bold).Add(color.Underline),
		//{38, 5, 156, 1, 4}, // user +x bit (files)
		"ue": FgColor256A(156).Add(color.Bold),
		//{38, 5, 156, 1},    // user +x bit (file types)
		"gr": FgColor256A(230).Add(color.Bold),
		//{38, 5, 230, 1}, // group +r bit
		"gw": FgColor256A(209).Add(color.Bold),
		//{38, 5, 209, 1}, // group +w bit
		"gx": FgColor256A(156).Add(color.Bold).Add(color.Underline),
		//{38, 5, 156, 1, 4}, // group +x bit
		"tr": FgColor256A(230).Add(color.Bold),
		//{38, 5, 230, 1}, // others +r bit
		"tw": FgColor256A(209).Add(color.Bold),
		//{38, 5, 209, 1}, // others +w bit
		"tx": FgColor256A(156).Add(color.Bold).Add(color.Underline),
		//{38, 5, 156, 1, 4}, // others +x bit
		"sn": FgColor256A(156).Add(color.Bold),
		//{38, 5, 156, 1}, // size number
		"snu": FgColor256A(156),
		//{38, 5, 156}, // size unit
		"uu": {38, 5, 229, 1},
		// user is you + 1 -> bold
		// "un": {38, 5, 214},    // user is not you
		"un": FgColor256A(GraysI[19]),
		//{38, 5, 251},    // user is not you
		"gu": {38, 5, 229, 1},
		// group with you in it
		// "gn": {38, 5, 214},    // group without you
		"gn": FgColor256A(GraysI[19]),
		//{38, 5, 251}, // group without you
		"da": FgColor256A(153),
		//{38, 5, 153}, // timestamp + 8 -> concealed
		"hd": NewAttributeA().Add(BgGrayA(3)...).Add(color.Underline),
		// "hd": FgGrayA(16).Add(BgGrayA(4)...).Add(color.Bold, color.Underline),
		// "hd": {4, 38, 5, 15}, // head
		// head + 4-> underline
		"-": FgColor256A(8),
		//{38, 5, 8}, // Concealed
		".": FgColor256A(8),
		//{38, 5, 8}, // Concealed
		" ": FgColor256A(8),
		//{38, 5, 8}, // Concealed
		"ga": FgColor256A(156),
		//{38, 5, 156}, // git new
		"gm": FgColor256A(39),
		//{38, 5, 39}, // git modified
		"gd": FgColor256A(196),
		//{38, 5, 196}, // git deleted
		"gv": FgColor256A(186),
		//{38, 5, 186}, // git renamed
		"gt": FgColor256A(207),
		//{38, 5, 207}, // git type change
		"dir": FgColor256A(189),
		//{38, 5, 189}, //addition 'dir'
		// "xattr": {38, 5, 249, 4}, //addition 'xattr'+ 4-> underline
		"xattr": FgColor256A(8).Add(BgGrayA(2)...).Add(color.Underline),
		//{38, 5, 8, color.Underline, 48, 5, Grays[2]},
		"xsymb": FgColor256A(8).Add(BgGrayA(2)...),
		//{38, 5, 8, 48, 5, Grays[2]},
		"in": FgColor256A(213),
		// {38, 5, 213}, // inode
		"lk": FgColor256A(209).Add(color.Bold),
		//{38, 5, 209, 1}, // links
		"bk": FgColor256A(189),
		// {38, 5, 189},                                              // blocks
		"pmpt": FgGrayA(19).Add(BgGrayA(4)...),
		// "prompt": FgColor256A(Grays[19]).Add(BgColor256A(Grays[4])...),
		//{38, 5, 251, 48, 5, 236},
		"bgpmpt": BgGrayA(4),
		//{48, 5, 236},
		"pmptsn": FgColor256A(156).Add(BgGrayA(4)...).Add(color.Bold),
		// "pmptsn": FgColor256A(156).Add(BgColor256A(Grays[4])...).Add(color.Bold),
		//{38, 5, 156, color.Bold, 48, 5, Grays[6]},
		"pmptsu": FgColor256A(156).Add(BgGrayA(4)...),
		// "pmptsu": FgColor256A(156).Add(BgColor256A(Grays[4])...),
		//{38, 5, 156, 48, 5, 236},
		"pmptdash": FgColor256A(8).Add(BgGrayA(4)...),
		"trace":    LogLevelColorA(logrus.TraceLevel),
		"debug":    LogLevelColorA(logrus.DebugLevel),
		"info":     LogLevelColorA(logrus.InfoLevel),
		"warn":     LogLevelColorA(logrus.WarnLevel),
		"error":    LogLevelColorA(logrus.ErrorLevel),
		"fatal":    LogLevelColorA(logrus.FatalLevel),
		"panic":    LogLevelColorA(logrus.PanicLevel),
		"md5":      LSColorAttributes["no"],
		//LSColorAttributes[".md5"],
		"field": FgColor256A(216),
		// {38, 5, 216},
		"value": FgColor256A(222).Add(color.Underline),
		//{38, 5, 222, 4},
		// Cvalue :{38, 5, 193, 4},
		"evenH": FgColor256A(231).Add(BgGrayA(4)...).Add(color.Underline),
		"even":  FgColor256A(231),
		// "even": FgGrayA(19).Add(BgGrayA(4)...),
		//{38, 5, 251, 48, 5, 236},
		"oddH": FgColor256A(223).Add(BgGrayA(4)...).Add(color.Underline),
		"odd":  FgColor256A(223),
		//{38, 5, 159, 48, 5, 238},
	}
	// LSColorAttributes = make(map[string]string) is LS_COLORS code according to
	// extention of file
	LSColorAttributes = map[string][]Attribute{
		"bd": FgColor256A(68),
		//{38, 5, 68},
		"ca": FgColor256A(17),
		//{38, 5, 17},
		"cd": FgColor256A(113).Add(color.Bold),
		//{38, 5, 113, color.Bold},
		"di": FgColor256A(30),
		//{38, 5, 30},
		"do": FgColor256A(127),
		//{38, 5, 127},
		"ex": FgColor256A(208).Add(color.Bold),
		//{38, 5, 208, color.Bold},
		"pi": FgColor256A(126),
		//{38, 5, 126},
		"fi": NewAttributeA(),
		//{0},
		"ln": FgColor256A(45),
		//{38, 5, 45},
		"mh": FgColor256A(222).Add(color.Bold),
		//{38, 5, 222, color.Bold},
		"no": FgColor256A(0),
		//{0},
		"or": FgColor256A(232).Add(BgColor256A(196)...).Add(color.Bold),
		//{48, 5, 196, 38, 5, 232, color.Bold},
		"ow": FgColor256A(220).Add(color.Bold),
		//{38, 5, 220, color.Bold},
		"sg": FgColor256A(0).Add(BgColor256A(0).Add(color.Italic)...),
		//{48, 5, color.Italic, 38, 5, 0},
		"su": FgColor256A(220).Add(color.Bold, color.Italic, color.BgHiBlack),
		//{38, 5, 220, color.Bold, color.Italic, color.BgHiBlack, color.Bold},
		"so": FgColor256A(197),
		//{38, 5, 197},
		"st": FgColor256A(86).Add(BgGrayA(2)...),
		//{38, 5, 86, 48, 5, 234},
		"tw": FgColor256A(235).Add(BgColor256A(139)...),
		//{48, 5, 235, 38, 5, 139, color.Italic},
		"LS_COLORS": FgColor256A(197).Add(BgColor256A(89)...).Add(color.Bold, color.Italic, color.Underline, color.ReverseVideo),
		//{48, 5, 89, 38, 5, 197, color.Bold, color.Italic, color.Underline, color.ReverseVideo},
		"-": FgColor256A(8),
		//{38, 5, 8}, // Concealed
		".": FgColor256A(8),
		//{38, 5, 8}, // Concealed
		"README": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"README.rst": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"README.md": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"LICENSE": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"COPYING": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"INSTALL": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"COPYRIGHT": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"AUTHORS": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"HISTORY": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"CONTRIBUTORS": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"PATENTS": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"VERSION": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"NOTICE": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		"CHANGES": FgColor256A(220).Add(color.Bold).Add(color.Underline),
		//{38, 5, 220, 1, 4},
		".log": FgColor256A(190),
		//{38, 5, 190},
		".txt": FgColor256A(253),
		//{38, 5, 253},
		".etx": FgColor256A(184),
		//{38, 5, 184},
		".info": FgColor256A(184),
		//{38, 5, 184},
		".markdown": FgColor256A(184),
		//{38, 5, 184},
		".md": FgColor256A(184),
		//{38, 5, 184},
		".mkd": FgColor256A(184),
		//{38, 5, 184},
		".nfo": FgColor256A(184),
		//{38, 5, 184},
		".pod": FgColor256A(184),
		//{38, 5, 184},
		".rst": FgColor256A(184),
		//{38, 5, 184},
		".tex": FgColor256A(184),
		//{38, 5, 184},
		".textile": FgColor256A(184),
		//{38, 5, 184},
		".bib": FgColor256A(178),
		//{38, 5, 178},
		".json": FgColor256A(178),
		//{38, 5, 178},
		".jsonl": FgColor256A(178),
		//{38, 5, 178},
		".jsonnet": FgColor256A(178),
		//{38, 5, 178},
		".libsonnet": FgColor256A(142),
		//{38, 5, 142},
		".ndjson": FgColor256A(178),
		//{38, 5, 178},
		".msg": FgColor256A(178),
		//{38, 5, 178},
		".pgn": FgColor256A(178),
		//{38, 5, 178},
		".rss": FgColor256A(178),
		//{38, 5, 178},
		".xml": FgColor256A(178),
		//{38, 5, 178},
		".fxml": FgColor256A(178),
		//{38, 5, 178},
		".toml": FgColor256A(178),
		//{38, 5, 178},
		".yaml": FgColor256A(178),
		//{38, 5, 178},
		".yml": FgColor256A(178),
		//{38, 5, 178},
		".RData": FgColor256A(178),
		//{38, 5, 178},
		".rdata": FgColor256A(178),
		//{38, 5, 178},
		".xsd": FgColor256A(178),
		//{38, 5, 178},
		".dtd": FgColor256A(178),
		//{38, 5, 178},
		".sgml": FgColor256A(178),
		//{38, 5, 178},
		".rng": FgColor256A(178),
		//{38, 5, 178},
		".rnc": FgColor256A(178),
		//{38, 5, 178},
		".cbr": FgColor256A(141),
		//{38, 5, 141},
		".cbz": FgColor256A(141),
		//{38, 5, 141},
		".chm": FgColor256A(141),
		//{38, 5, 141},
		".djvu": FgColor256A(141),
		//{38, 5, 141},
		".pdf": FgColor256A(141),
		//{38, 5, 141},
		".PDF": FgColor256A(141),
		//{38, 5, 141},
		".mobi": FgColor256A(141),
		//{38, 5, 141},
		".epub": FgColor256A(141),
		//{38, 5, 141},
		".docm": FgColor256A(111).Add(color.Underline),
		//{38, 5, 111, color.Underline},
		".doc": FgColor256A(111),
		//{38, 5, 111},
		".docx": FgColor256A(111),
		//{38, 5, 111},
		".odb": FgColor256A(111),
		//{38, 5, 111},
		".odt": FgColor256A(111),
		//{38, 5, 111},
		".rtf": FgColor256A(111),
		//{38, 5, 111},
		".odp": FgColor256A(166),
		//{38, 5, 166},
		".pps": FgColor256A(166),
		//{38, 5, 166},
		".ppt": FgColor256A(166),
		//{38, 5, 166},
		".pptx": FgColor256A(166),
		//{38, 5, 166},
		".ppts": FgColor256A(166),
		//{38, 5, 166},
		".pptxm": FgColor256A(166).Add(color.Underline),
		//{38, 5, 166, color.Underline},
		".pptsm": FgColor256A(166).Add(color.Underline),
		//{38, 5, 166, color.Underline},
		".csv": FgColor256A(78),
		//{38, 5, 78},
		".tsv": FgColor256A(78),
		//{38, 5, 78},
		".ods": FgColor256A(112),
		//{38, 5, 112},
		".xla": FgColor256A(76),
		//{38, 5, 76},
		".xls": FgColor256A(112),
		//{38, 5, 112},
		".xlsx": FgColor256A(112),
		//{38, 5, 112},
		".xlsxm": FgColor256A(112).Add(color.Underline),
		//{38, 5, 112, color.Underline},
		".xltm": FgColor256A(73).Add(color.Underline),
		//{38, 5, 73, color.Underline},
		".xltx": FgColor256A(73),
		//{38, 5, 73},
		".pages": FgColor256A(111),
		//{38, 5, 111},
		".numbers": FgColor256A(112),
		//{38, 5, 112},
		".key": FgColor256A(166),
		//{38, 5, 166},
		"config": NewAttributeA().Add(color.Bold),

		"cfg":                  NewAttributeA().Add(color.Bold),
		"conf":                 NewAttributeA().Add(color.Bold),
		"rc":                   NewAttributeA().Add(color.Bold),
		"authorized_keys":      NewAttributeA().Add(color.Bold),
		"known_hosts":          NewAttributeA().Add(color.Bold),
		".ini":                 NewAttributeA().Add(color.Bold),
		".plist":               NewAttributeA().Add(color.Bold),
		".viminfo":             NewAttributeA().Add(color.Bold),
		".pcf":                 NewAttributeA().Add(color.Bold),
		".psf":                 NewAttributeA().Add(color.Bold),
		".hidden-color-scheme": NewAttributeA().Add(color.Bold),
		".hidden-tmTheme":      NewAttributeA().Add(color.Bold),
		".last-run":            NewAttributeA().Add(color.Bold),
		".merged-ca-bundle":    NewAttributeA().Add(color.Bold),
		".sublime-build":       NewAttributeA().Add(color.Bold),
		".sublime-commands":    NewAttributeA().Add(color.Bold),
		".sublime-keymap":      NewAttributeA().Add(color.Bold),
		".sublime-settings":    NewAttributeA().Add(color.Bold),
		".sublime-snippet":     NewAttributeA().Add(color.Bold),
		".sublime-project":     NewAttributeA().Add(color.Bold),
		".sublime-workspace":   NewAttributeA().Add(color.Bold),
		".tmTheme":             NewAttributeA().Add(color.Bold),
		".user-ca-bundle":      NewAttributeA().Add(color.Bold),
		".epf":                 NewAttributeA().Add(color.Bold),
		".git":                 FgColor256A(197),
		//{38, 5, 197},
		".gitignore": FgColor256A(240),
		//{38, 5, 240},
		".gitattributes": FgColor256A(240),
		//{38, 5, 240},
		".gitmodules": FgColor256A(240),
		//{38, 5, 240},
		".awk": FgColor256A(172),
		//{38, 5, 172},
		".bash": FgColor256A(172),
		//{38, 5, 172},
		".bat": FgColor256A(172),
		//{38, 5, 172},
		".BAT": FgColor256A(172),
		//{38, 5, 172},
		".sed": FgColor256A(172),
		//{38, 5, 172},
		".sh": FgColor256A(172),
		//{38, 5, 172},
		".zsh": FgColor256A(172),
		//{38, 5, 172},
		".vim": FgColor256A(172),
		//{38, 5, 172},
		".kak": FgColor256A(172),
		//{38, 5, 172},
		".ahk": FgColor256A(41),
		//{38, 5, 41},
		".py": FgColor256A(41),
		//{38, 5, 41},
		".ipynb": FgColor256A(41),
		//{38, 5, 41},
		".rb": FgColor256A(41),
		//{38, 5, 41},
		".gemspec": FgColor256A(41),
		//{38, 5, 41},
		".pl": FgColor256A(208),
		//{38, 5, 208},
		".PL": FgColor256A(160),
		//{38, 5, 160},
		".t": FgColor256A(114),
		//{38, 5, 114},
		".msql": FgColor256A(222),
		//{38, 5, 222},
		".mysql": FgColor256A(222),
		//{38, 5, 222},
		".pgsql": FgColor256A(222),
		//{38, 5, 222},
		".sql": FgColor256A(222),
		//{38, 5, 222},
		".tcl": FgColor256A(64).Add(color.Bold),
		//{38, 5, 64, color.Bold},
		".r": FgColor256A(49),
		//{38, 5, 49},
		".R": FgColor256A(49),
		//{38, 5, 49},
		".gs": FgColor256A(81),
		//{38, 5, 81},
		".clj": FgColor256A(41),
		//{38, 5, 41},
		".cljs": FgColor256A(41),
		//{38, 5, 41},
		".cljc": FgColor256A(41),
		//{38, 5, 41},
		".cljw": FgColor256A(41),
		//{38, 5, 41},
		".scala": FgColor256A(41),
		//{38, 5, 41},
		".sc": FgColor256A(41),
		//{38, 5, 41},
		".dart": FgColor256A(51),
		//{38, 5, 51},
		".asm": FgColor256A(81),
		//{38, 5, 81},
		".cl": FgColor256A(81),
		//{38, 5, 81},
		".lisp": FgColor256A(81),
		//{38, 5, 81},
		".rkt": FgColor256A(81),
		//{38, 5, 81},
		".lua": FgColor256A(81),
		//{38, 5, 81},
		".moon": FgColor256A(81),
		//{38, 5, 81},
		".c": FgColor256A(81),
		//{38, 5, 81},
		".C": FgColor256A(81),
		//{38, 5, 81},
		".h": FgColor256A(110),
		//{38, 5, 110},
		".H": FgColor256A(110),
		//{38, 5, 110},
		".tcc": FgColor256A(110),
		//{38, 5, 110},
		".c++": FgColor256A(81),
		//{38, 5, 81},
		".h++": FgColor256A(110),
		//{38, 5, 110},
		".hpp": FgColor256A(110),
		//{38, 5, 110},
		".hxx": FgColor256A(110),
		//{38, 5, 110},
		".ii": FgColor256A(110),
		//{38, 5, 110},
		".M": FgColor256A(110),
		//{38, 5, 110},
		".m": FgColor256A(110),
		//{38, 5, 110},
		".cc": FgColor256A(81),
		//{38, 5, 81},
		".cs": FgColor256A(81),
		//{38, 5, 81},
		".cp": FgColor256A(81),
		//{38, 5, 81},
		".cpp": FgColor256A(81),
		//{38, 5, 81},
		".cxx": FgColor256A(81),
		//{38, 5, 81},
		".cr": FgColor256A(81),
		//{38, 5, 81},
		".go": FgColor256A(81),
		//{38, 5, 81},
		".f": FgColor256A(81),
		//{38, 5, 81},
		".F": FgColor256A(81),
		//{38, 5, 81},
		".for": FgColor256A(81),
		//{38, 5, 81},
		".ftn": FgColor256A(81),
		//{38, 5, 81},
		".f90": FgColor256A(81),
		//{38, 5, 81},
		".F90": FgColor256A(81),
		//{38, 5, 81},
		".f95": FgColor256A(81),
		//{38, 5, 81},
		".F95": FgColor256A(81),
		//{38, 5, 81},
		".f03": FgColor256A(81),
		//{38, 5, 81},
		".F03": FgColor256A(81),
		//{38, 5, 81},
		".f08": FgColor256A(81),
		//{38, 5, 81},
		".F08": FgColor256A(81),
		//{38, 5, 81},
		".nim": FgColor256A(81),
		//{38, 5, 81},
		".nimble": FgColor256A(81),
		//{38, 5, 81},
		".s": FgColor256A(110),
		//{38, 5, 110},
		".S": FgColor256A(110),
		//{38, 5, 110},
		".rs": FgColor256A(81),
		//{38, 5, 81},
		".scpt": FgColor256A(219),
		//{38, 5, 219},
		".swift": FgColor256A(219),
		//{38, 5, 219},
		".sx": FgColor256A(81),
		//{38, 5, 81},
		".vala": FgColor256A(81),
		//{38, 5, 81},
		".vapi": FgColor256A(81),
		//{38, 5, 81},
		".hi": FgColor256A(110),
		//{38, 5, 110},
		".hs": FgColor256A(81),
		//{38, 5, 81},
		".lhs": FgColor256A(81),
		//{38, 5, 81},
		".agda": FgColor256A(81),
		//{38, 5, 81},
		".lagda": FgColor256A(81),
		//{38, 5, 81},
		".lagda.tex": FgColor256A(81),
		//{38, 5, 81},
		".lagda.rst": FgColor256A(81),
		//{38, 5, 81},
		".lagda.md": FgColor256A(81),
		//{38, 5, 81},
		".agdai": FgColor256A(110),
		//{38, 5, 110},
		".zig": FgColor256A(81),
		//{38, 5, 81},
		".v": FgColor256A(81),
		//{38, 5, 81},
		".pyc": FgColor256A(240),
		//{38, 5, 240},
		".tf": FgColor256A(168),
		//{38, 5, 168},
		".tfstate": FgColor256A(168),
		//{38, 5, 168},
		".tfvars": FgColor256A(168),
		//{38, 5, 168},
		".css": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".less": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".sass": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".scss": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".htm": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".html": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".jhtm": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".mht": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".eml": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".mustache": FgColor256A(125).Add(color.Bold),
		//{38, 5, 125, color.Bold},
		".coffee": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".java": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".js": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".mjs": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".jsm": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".jsp": FgColor256A(74).Add(color.Bold),
		//{38, 5, 074, color.Bold},
		".php": FgColor256A(81),
		//{38, 5, 81},
		".ctp": FgColor256A(81),
		//{38, 5, 81},
		".twig": FgColor256A(81),
		//{38, 5, 81},
		".vb": FgColor256A(81),
		//{38, 5, 81},
		".vba": FgColor256A(81),
		//{38, 5, 81},
		".vbs": FgColor256A(81),
		//{38, 5, 81},
		"Dockerfile": FgColor256A(155).Add(color.Underline),
		//{38, 5, 155, color.Underline},
		".dockerignore": FgColor256A(240),
		//{38, 5, 240},
		"Makefile": FgColor256A(155).Add(color.Underline),
		//{38, 5, 155, color.Underline},
		"MANIFEST": FgColor256A(243).Add(color.Underline),
		//{38, 5, 243, color.Underline},
		"pm_to_blib": FgColor256A(240),
		//{38, 5, 240},
		".nix": FgColor256A(155),
		//{38, 5, 155},
		".dhall": FgColor256A(178),
		//{38, 5, 178},
		".rake": FgColor256A(155),
		//{38, 5, 155},
		".am": FgColor256A(242),
		//{38, 5, 242},
		".in": FgColor256A(242),
		//{38, 5, 242},
		".hin": FgColor256A(242),
		//{38, 5, 242},
		".scan": FgColor256A(242),
		//{38, 5, 242},
		".m4": FgColor256A(242),
		//{38, 5, 242},
		".old": FgColor256A(242),
		//{38, 5, 242},
		".out": FgColor256A(242),
		//{38, 5, 242},
		".SKIP": FgColor256A(244),
		//{38, 5, 244},
		".diff": FgColor256A(197).Add(BgColor256A(232)...),
		//{48, 5, 197, 38, 5, 232},
		".patch": FgColor256A(197).Add(BgColor256A(232)...).Add(color.Bold),
		//{48, 5, 197, 38, 5, 232, color.Bold},
		".bmp": FgColor256A(97),
		//{38, 5, 97},
		".dicom": FgColor256A(97),
		//{38, 5, 97},
		".tiff": FgColor256A(97),
		//{38, 5, 97},
		".tif": FgColor256A(97),
		//{38, 5, 97},
		".TIFF": FgColor256A(97),
		//{38, 5, 97},
		".cdr": FgColor256A(97),
		//{38, 5, 97},
		".flif": FgColor256A(97),
		//{38, 5, 97},
		".gif": FgColor256A(97),
		//{38, 5, 97},
		".icns": FgColor256A(97),
		//{38, 5, 97},
		".ico": FgColor256A(97),
		//{38, 5, 97},
		".jpeg": FgColor256A(97),
		//{38, 5, 97},
		".JPG": FgColor256A(97),
		//{38, 5, 97},
		".jpg": FgColor256A(97),
		//{38, 5, 97},
		".nth": FgColor256A(97),
		//{38, 5, 97},
		".png": FgColor256A(97),
		//{38, 5, 97},
		".psd": FgColor256A(97),
		//{38, 5, 97},
		".pxd": FgColor256A(97),
		//{38, 5, 97},
		".pxm": FgColor256A(97),
		//{38, 5, 97},
		".xpm": FgColor256A(97),
		//{38, 5, 97},
		".webp": FgColor256A(97),
		//{38, 5, 97},
		".ai": FgColor256A(99),
		//{38, 5, 99},
		".eps": FgColor256A(99),
		//{38, 5, 99},
		".epsf": FgColor256A(99),
		//{38, 5, 99},
		".drw": FgColor256A(99),
		//{38, 5, 99},
		".ps": FgColor256A(99),
		//{38, 5, 99},
		".svg": FgColor256A(99),
		//{38, 5, 99},
		".avi": FgColor256A(114),
		//{38, 5, 114},
		".divx": FgColor256A(114),
		//{38, 5, 114},
		".IFO": FgColor256A(114),
		//{38, 5, 114},
		".m2v": FgColor256A(114),
		//{38, 5, 114},
		".m4v": FgColor256A(114),
		//{38, 5, 114},
		".mkv": FgColor256A(114),
		//{38, 5, 114},
		".MOV": FgColor256A(114),
		//{38, 5, 114},
		".mov": FgColor256A(114),
		//{38, 5, 114},
		".mp4": FgColor256A(114),
		//{38, 5, 114},
		".mpeg": FgColor256A(114),
		//{38, 5, 114},
		".mpg": FgColor256A(114),
		//{38, 5, 114},
		".ogm": FgColor256A(114),
		//{38, 5, 114},
		".rmvb": FgColor256A(114),
		//{38, 5, 114},
		".sample": FgColor256A(114),
		//{38, 5, 114},
		".wmv": FgColor256A(114),
		//{38, 5, 114},
		".3g2": FgColor256A(115),
		//{38, 5, 115},
		".3gp": FgColor256A(115),
		//{38, 5, 115},
		".gp3": FgColor256A(115),
		//{38, 5, 115},
		".webm": FgColor256A(115),
		//{38, 5, 115},
		".gp4": FgColor256A(115),
		//{38, 5, 115},
		".asf": FgColor256A(115),
		//{38, 5, 115},
		".flv": FgColor256A(115),
		//{38, 5, 115},
		".ts": FgColor256A(115),
		//{38, 5, 115},
		".ogv": FgColor256A(115),
		//{38, 5, 115},
		".f4v": FgColor256A(115),
		//{38, 5, 115},
		".VOB": FgColor256A(115).Add(color.Bold),
		//{38, 5, 115, color.Bold},
		".vob": FgColor256A(115).Add(color.Bold),
		//{38, 5, 115, color.Bold},
		".ass": FgColor256A(117),
		//{38, 5, 117},
		".srt": FgColor256A(117),
		//{38, 5, 117},
		".ssa": FgColor256A(117),
		//{38, 5, 117},
		".sub": FgColor256A(117),
		//{38, 5, 117},
		".sup": FgColor256A(117),
		//{38, 5, 117},
		".vtt": FgColor256A(117),
		//{38, 5, 117},
		".3ga": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".S3M": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".aac": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".amr": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".au": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".caf": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".dat": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".dts": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".fcm": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".m4a": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".mod": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".mp3": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".mp4a": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".oga": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".ogg": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".opus": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".s3m": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".sid": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".wma": FgColor256A(137).Add(color.Bold),
		//{38, 5, 137, color.Bold},
		".ape": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".aiff": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".cda": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".flac": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".alac": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".mid": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".midi": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".pcm": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".wav": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".wv": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".wvc": FgColor256A(136).Add(color.Bold),
		//{38, 5, 136, color.Bold},
		".afm": FgColor256A(66),
		//{38, 5, 66},
		".fon": FgColor256A(66),
		//{38, 5, 66},
		".fnt": FgColor256A(66),
		//{38, 5, 66},
		".pfb": FgColor256A(66),
		//{38, 5, 66},
		".pfm": FgColor256A(66),
		//{38, 5, 66},
		".ttf": FgColor256A(66),
		//{38, 5, 66},
		".otf": FgColor256A(66),
		//{38, 5, 66},
		".woff": FgColor256A(66),
		//{38, 5, 66},
		".woff2": FgColor256A(66),
		//{38, 5, 66},
		".PFA": FgColor256A(66),
		//{38, 5, 66},
		".pfa": FgColor256A(66),
		//{38, 5, 66},
		".7z": FgColor256A(40),
		//{38, 5, 40},
		".a": FgColor256A(40),
		//{38, 5, 40},
		".arj": FgColor256A(40),
		//{38, 5, 40},
		".bz2": FgColor256A(40),
		//{38, 5, 40},
		".cpio": FgColor256A(40),
		//{38, 5, 40},
		".gz": FgColor256A(40),
		//{38, 5, 40},
		".lrz": FgColor256A(40),
		//{38, 5, 40},
		".lz": FgColor256A(40),
		//{38, 5, 40},
		".lzma": FgColor256A(40),
		//{38, 5, 40},
		".lzo": FgColor256A(40),
		//{38, 5, 40},
		".rar": FgColor256A(40),
		//{38, 5, 40},
		".s7z": FgColor256A(40),
		//{38, 5, 40},
		".sz": FgColor256A(40),
		//{38, 5, 40},
		".tar": FgColor256A(40),
		//{38, 5, 40},
		".tgz": FgColor256A(40),
		//{38, 5, 40},
		".warc": FgColor256A(40),
		//{38, 5, 40},
		".WARC": FgColor256A(40),
		//{38, 5, 40},
		".xz": FgColor256A(40),
		//{38, 5, 40},
		".z": FgColor256A(40),
		//{38, 5, 40},
		".zip": FgColor256A(40),
		//{38, 5, 40},
		".zipx": FgColor256A(40),
		//{38, 5, 40},
		".zoo": FgColor256A(40),
		//{38, 5, 40},
		".zpaq": FgColor256A(40),
		//{38, 5, 40},
		".zst": FgColor256A(40),
		//{38, 5, 40},
		".zstd": FgColor256A(40),
		//{38, 5, 40},
		".zz": FgColor256A(40),
		//{38, 5, 40},
		".apk": FgColor256A(215),
		//{38, 5, 215},
		".ipa": FgColor256A(215),
		//{38, 5, 215},
		".deb": FgColor256A(215),
		//{38, 5, 215},
		".rpm": FgColor256A(215),
		//{38, 5, 215},
		".jad": FgColor256A(215),
		//{38, 5, 215},
		".jar": FgColor256A(215),
		//{38, 5, 215},
		".cab": FgColor256A(215),
		//{38, 5, 215},
		".pak": FgColor256A(215),
		//{38, 5, 215},
		".pk3": FgColor256A(215),
		//{38, 5, 215},
		".vdf": FgColor256A(215),
		//{38, 5, 215},
		".vpk": FgColor256A(215),
		//{38, 5, 215},
		".bsp": FgColor256A(215),
		//{38, 5, 215},
		".dmg": FgColor256A(215),
		//{38, 5, 215},
		".r[0-9]{0,2}": FgColor256A(239),
		//{38, 5, 239},
		".zx[0-9]{0,2}": FgColor256A(239),
		//{38, 5, 239},
		".z[0-9]{0,2}": FgColor256A(239),
		//{38, 5, 239},
		".part": FgColor256A(239),
		//{38, 5, 239},
		".iso": FgColor256A(124),
		//{38, 5, 124},
		".bin": FgColor256A(124),
		//{38, 5, 124},
		".nrg": FgColor256A(124),
		//{38, 5, 124},
		".qcow": FgColor256A(124),
		//{38, 5, 124},
		".sparseimage": FgColor256A(124),
		//{38, 5, 124},
		".toast": FgColor256A(124),
		//{38, 5, 124},
		".vcd": FgColor256A(124),
		//{38, 5, 124},
		".vmdk": FgColor256A(124),
		//{38, 5, 124},
		".accdb": FgColor256A(60),
		//{38, 5, 60},
		".accde": FgColor256A(60),
		//{38, 5, 60},
		".accdr": FgColor256A(60),
		//{38, 5, 60},
		".accdt": FgColor256A(60),
		//{38, 5, 60},
		".db": FgColor256A(60),
		//{38, 5, 60},
		".fmp12": FgColor256A(60),
		//{38, 5, 60},
		".fp7": FgColor256A(60),
		//{38, 5, 60},
		".localstorage": FgColor256A(60),
		//{38, 5, 60},
		".mdb": FgColor256A(60),
		//{38, 5, 60},
		".mde": FgColor256A(60),
		//{38, 5, 60},
		".sqlite": FgColor256A(60),
		//{38, 5, 60},
		".typelib": FgColor256A(60),
		//{38, 5, 60},
		".nc": FgColor256A(60),
		//{38, 5, 60},
		".pacnew": FgColor256A(33),
		//{38, 5, 33},
		".un~": FgColor256A(241),
		//{38, 5, 241},
		".orig": FgColor256A(241),
		//{38, 5, 241},
		".BUP": FgColor256A(241),
		//{38, 5, 241},
		".bak": FgColor256A(241),
		//{38, 5, 241},
		".o": FgColor256A(241),
		//{38, 5, 241},
		"core": FgColor256A(241),
		//{38, 5, 241},
		".mdump": FgColor256A(241),
		//{38, 5, 241},
		".rlib": FgColor256A(241),
		//{38, 5, 241},
		".dll": FgColor256A(241),
		//{38, 5, 241},
		".swp": FgColor256A(244),
		//{38, 5, 244},
		".swo": FgColor256A(244),
		//{38, 5, 244},
		".tmp": FgColor256A(244),
		//{38, 5, 244},
		".sassc": FgColor256A(244),
		//{38, 5, 244},
		".pid": FgColor256A(248),
		//{38, 5, 248},
		".state": FgColor256A(248),
		//{38, 5, 248},
		"lockfile": FgColor256A(248),
		//{38, 5, 248},
		"lock": FgColor256A(248),
		//{38, 5, 248},
		".err": FgColor256A(160).Add(color.Bold),
		//{38, 5, 160, color.Bold},
		".error": FgColor256A(160).Add(color.Bold),
		//{38, 5, 160, color.Bold},
		".stderr": FgColor256A(160).Add(color.Bold),
		//{38, 5, 160, color.Bold},
		".aria2": FgColor256A(241),
		//{38, 5, 241},
		".dump": FgColor256A(241),
		//{38, 5, 241},
		".stackdump": FgColor256A(241),
		//{38, 5, 241},
		".zcompdump": FgColor256A(241),
		//{38, 5, 241},
		".zwc": FgColor256A(241),
		//{38, 5, 241},
		".pcap": FgColor256A(29),
		//{38, 5, 29},
		".cap": FgColor256A(29),
		//{38, 5, 29},
		".dmp": FgColor256A(29),
		//{38, 5, 29},
		".DS_Store": FgColor256A(239),
		//{38, 5, 239},
		".localized": FgColor256A(239),
		//{38, 5, 239},
		".CFUserTextEncoding": FgColor256A(239),
		//{38, 5, 239},
		".allow": FgColor256A(112),
		//{38, 5, 112},
		".deny": FgColor256A(196),
		//{38, 5, 196},
		".service": FgColor256A(45),
		//{38, 5, 45},
		"@.service": FgColor256A(45),
		//{38, 5, 45},
		".socket": FgColor256A(45),
		//{38, 5, 45},
		".swap": FgColor256A(45),
		//{38, 5, 45},
		".device": FgColor256A(45),
		//{38, 5, 45},
		".mount": FgColor256A(45),
		//{38, 5, 45},
		".automount": FgColor256A(45),
		//{38, 5, 45},
		".target": FgColor256A(45),
		//{38, 5, 45},
		".path": FgColor256A(45),
		//{38, 5, 45},
		".timer": FgColor256A(45),
		//{38, 5, 45},
		".snapshot": FgColor256A(45),
		//{38, 5, 45},
		".application": FgColor256A(116),
		//{38, 5, 116},
		".cue": FgColor256A(116),
		//{38, 5, 116},
		".description": FgColor256A(116),
		//{38, 5, 116},
		".directory": FgColor256A(116),
		//{38, 5, 116},
		".m3u": FgColor256A(116),
		//{38, 5, 116},
		".m3u8": FgColor256A(116),
		//{38, 5, 116},
		".md5": FgColor256A(116),
		//{38, 5, 116},
		".properties": FgColor256A(116),
		//{38, 5, 116},
		".sfv": FgColor256A(116),
		//{38, 5, 116},
		".theme": FgColor256A(116),
		//{38, 5, 116},
		".torrent": FgColor256A(116),
		//{38, 5, 116},
		".urlview": FgColor256A(116),
		//{38, 5, 116},
		".webloc": FgColor256A(116),
		//{38, 5, 116},
		".lnk": FgColor256A(39),
		//{38, 5, 39},
		"CodeResources": FgColor256A(239),
		//{38, 5, 239},
		"PkgInfo": FgColor256A(239),
		//{38, 5, 239},
		".nib": FgColor256A(57),
		//{38, 5, 57},
		".car": FgColor256A(57),
		//{38, 5, 57},
		".dylib": FgColor256A(241),
		//{38, 5, 241},
		".entitlements": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".pbxproj": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".strings": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".storyboard": FgColor256A(196),
		//{38, 5, 196},
		".xcconfig": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".xcsettings": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".xcuserstate": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".xcworkspacedata": NewAttributeA().Add(color.Bold),
		//{color.Bold},
		".xib": FgColor256A(208),
		//{38, 5, 208},
		".asc": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".bfe": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".enc": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".gpg": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".signature": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".sig": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".p12": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".pem": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".pgp": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".p7s": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		"id_dsa": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		"id_rsa": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		"id_ecdsa": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		"id_ed25519": FgColor256A(192).Add(color.Italic),
		//{38, 5, 192, color.Italic},
		".32x": FgColor256A(213),
		// {38, 5, 213},
		".cdi": FgColor256A(213),
		// {38, 5, 213},
		".fm2": FgColor256A(213),
		// {38, 5, 213},
		".rom": FgColor256A(213),
		// {38, 5, 213},
		".sav": FgColor256A(213),
		// {38, 5, 213},
		".st": FgColor256A(213),
		// {38, 5, 213},
		".a00": FgColor256A(213),
		// {38, 5, 213},
		".a52": FgColor256A(213),
		// {38, 5, 213},
		".A64": FgColor256A(213),
		// {38, 5, 213},
		".a64": FgColor256A(213),
		// {38, 5, 213},
		".a78": FgColor256A(213),
		// {38, 5, 213},
		".adf": FgColor256A(213),
		// {38, 5, 213},
		".atr": FgColor256A(213),
		// {38, 5, 213},
		".gb": FgColor256A(213),
		// {38, 5, 213},
		".gba": FgColor256A(213),
		// {38, 5, 213},
		".gbc": FgColor256A(213),
		// {38, 5, 213},
		".gel": FgColor256A(213),
		// {38, 5, 213},
		".gg": FgColor256A(213),
		// {38, 5, 213},
		".ggl": FgColor256A(213),
		// {38, 5, 213},
		".ipk": FgColor256A(213),
		// {38, 5, 213},
		".j64": FgColor256A(213),
		// {38, 5, 213},
		".nds": FgColor256A(213),
		// {38, 5, 213},
		".nes": FgColor256A(213),
		// {38, 5, 213},
		".sms": FgColor256A(213),
		// {38, 5, 213},
		".8xp": FgColor256A(121),
		// {38, 5, 121},
		".8eu": FgColor256A(121),
		// {38, 5, 121},
		".82p": FgColor256A(121),
		// {38, 5, 121},
		".83p": FgColor256A(121),
		// {38, 5, 121},
		".8xe": FgColor256A(121),
		// {38, 5, 121},
		".stl": FgColor256A(216),
		// {38, 5, 216},
		".dwg": FgColor256A(216),
		// {38, 5, 216},
		".ply": FgColor256A(216),
		// {38, 5, 216},
		".wrl": FgColor256A(216),
		// {38, 5, 216},
		".pot": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".pcb": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".mm": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".gbr": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".scm": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".xcf": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".spl": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".Rproj": FgColor256A(11),
		//{38, 5, 11},
		".sis": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".1p": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".3p": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".cnc": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".def": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".ex": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".example": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".feature": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".ger": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".ics": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".map": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".mf": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".mfasl": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".mi": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".mtx": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".pc": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".pi": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".plt": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".pm": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".rdf": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".ru": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".sch": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".sty": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".sug": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".tdy": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".tfm": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".tfnt": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".tg": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".vcard": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".vcf": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".xln": NewAttributeA().Add(color.ReverseVideo),
		//{38, 5, color.ReverseVideo},
		".iml": FgColor256A(166),
		//{38, 5, 166},
	}
	// ReExtLSColors is LS_COLORS code for specific pattern of file extentions
	ReExtLSColors = map[*regexp.Regexp][]Attribute{
		regexp.MustCompile(`r[0-9]{0,2}$`): FgColor256A(GraysI[7]),
		//{38, 5, 239},
		regexp.MustCompile(`zx[0-9]{0,2}$`): FgColor256A(GraysI[7]),
		//{38, 5, 239},
		regexp.MustCompile(`z[0-9]{0,2}$`): FgColor256A(GraysI[7]),
		//{38, 5, 239},
	}

	// Chdp is default color use for head
	Chdp = NewEXAColor("hd")
	// Cdirp is default color use for dir part of path
	Cdirp = NewEXAColor("dir")
	// Cdip is default color use for directory
	Cdip = NewEXAColor("di") //
	// Cfip is default color use for file
	Cfip = NewEXAColor("fi")
	// CNop is default color use for serial number
	CNop = NewEXAColor("-") //
	// Cinp is default color use for inode field
	Cinp = NewEXAColor("in")
	// Cpmp is default color use for permission field
	Cpms = NewEXAColor("uw")
	// Csnp is default color use for number of size
	Csnp = NewEXAColor("sn")
	// Csup is default color use for unit of size
	Csup = NewEXAColor("snu")
	// Cuup is default color use for user field
	Cuup = NewEXAColor("uu")
	// Cgup is default color use for group field
	Cgup = NewEXAColor("gu")
	// Cunp is default color use for user field, but user is not you
	Cunp = NewEXAColor("un")
	// Cgnp is default color use for group field, but group without you
	Cgnp = NewEXAColor("gn")
	// Clkp is default color use for hard link field
	Clkp = NewEXAColor("lk")
	// Cbkp is default color use for blocks field
	Cbkp = NewEXAColor("bk")
	// Cdap is default color use for date field
	Cdap = NewEXAColor("da")
	// Cgitp is default color use for git field
	Cgitp = NewEXAColor("gm")
	// Cmd5p is default color use for md5 field
	Cmd5p = NewEXAColor("md5")
	// Cxap is default color use for extended attributes
	Cxap = NewEXAColor("xattr")
	// Cxbp is default color use for symbole of extended attributes
	Cxbp = NewEXAColor("xsymb")
	// Cdashp is default color use for dash
	Cdashp = NewEXAColor("-")
	// Cnop is default color use for no this file kind
	Cnop = NewEXAColor("no")
	// Cbdp is default color use for device
	Cbdp = NewLSColor("bd")
	// Cbdp is default color use for chardevice
	Ccdp = NewLSColor("cd")
	// Cpip is default color use for named pipe (FIFO)
	Cpip = NewLSColor("pi")
	// Csop is default color use for socket
	Csop = NewLSColor("so")
	// Clnp is default color use for symlink
	Clnp = NewEXAColor("ln")
	// Cexp is default color use for execution file (permission contains 'x')
	Cexp = NewLSColor("ex")
	// Corp is default color use for orphan file. "orphan, symbolic link pointing to a non-existent file (orphan)"
	Corp = NewLSColor("or")
	// Cprompt is default color use for prompt
	Cpmpt = NewEXAColor("pmpt")
	// CpmptSn is default color use for number in prompt
	CpmptSn = NewEXAColor("pmptsn")
	// CpmptSu is default color use for unit in prompt
	CpmptSu    = NewEXAColor("pmptsu")
	CpmptDashp = NewEXAColor("pmptdash")
	Ctrace     = NewEXAColor("trace")
	Cdebug     = NewEXAColor("debug")
	Cinfo      = NewEXAColor("info")
	Cwarn      = NewEXAColor("warn")
	Cerror     = NewEXAColor("error")
	Cfatal     = NewEXAColor("fatal")
	Cpanic     = NewEXAColor("panic")
	Cfield     = NewEXAColor("field")
	Cvalue     = NewEXAColor("value")
	CEvenH     = NewEXAColor("evenH")
	COddH      = NewEXAColor("oddH")
	CEven      = NewEXAColor("even")
	COdd       = NewEXAColor("odd")
)

func ChoseColors(i int, colors []*Color) *Color {
	n := len(colors)
	if n == 0 {
		return Cnop
	}
	return colors[i%n]
}

var (
	_c2s  = []*Color{CEven, COdd}
	_c2Hs = []*Color{CEvenH, COddH}
)

func ChoseColor(i int) *Color {
	return _c2s[i%2]
}

func ChoseColorH(i int) *Color {
	return _c2Hs[i%2]
}

func CloneColor(color *Color) *Color {
	c := *color
	return &c
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var reANSI = regexp.MustCompile(ansi)

// StripANSI returns a string without ESC color code
func StripANSI(str string) string {
	return reANSI.ReplaceAllString(str, "")
}

// DisableColor will set `true` to `NoColor`
func DisableColor() {
	NoColor = true
	color.NoColor = NoColor
}

// EnableColor will set `true` to `NoColor`
func EnableColor() {
	NoColor = false
	color.NoColor = NoColor
}

// DefaultNoColor will resume the default value of `NoColor`
func DefaultNoColor() {
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
	color.NoColor = NoColor
}

// GetColors get LS_COLORS from env of os
func GetLSColors() {
	colorenv := os.Getenv("LS_COLORS")
	args := strings.Split(colorenv, ":")

	for _, a := range args {
		kv := strings.Split(a, "=")
		if len(kv) == 2 {
			LSColorAttributes[kv[0]] = getLSColorAttribute(kv[1])
		}
	}
}

func getLSColorAttribute(code string) []Attribute {
	att := []Attribute{}
	for _, a := range strings.Split(code, ";") {
		att = append(att, Attribute(cast.ToInt(a)))
	}
	return att
}

// KindLSColorString will colorful string `s` using key `kind`
func KindLSColorString(kind, s string) string {
	att, ok := LSColorAttributes[kind]
	if !ok {
		att = LSColorAttributes["fi"]
	}
	return color.New(att...).Sprint(s)
}

func KindEXAColorString(kind, s string) string {
	att, ok := EXAColorAttributes[kind]
	if !ok {
		att = EXAColorAttributes["fi"]
	}
	return color.New(att...).Sprint(s)
}

func LogLevelColorA(level logrus.Level) (a []Attribute) {
	switch level {
	case logrus.TraceLevel:
		a = []Attribute{color.FgCyan}
	case logrus.DebugLevel:
		a = []Attribute{color.FgHiRed}
	case logrus.WarnLevel:
		a = []Attribute{color.FgHiRed}
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		// a = FgColor256A(220).Add(color.Bold).Add(BgColor256A(160)...)
		a = []Attribute{38, 5, 220, 1, 48, 5, 160}
	default: //info
		a = []Attribute{color.FgHiGreen}
	}
	return a
}
func LogLevelColor(level logrus.Level) (c *Color) {
	a := LogLevelColorA(level)
	return color.New(a...)
}

// NewLSColor will return `*color.Color` using `LSColorAttributes[key]`
func NewLSColor(key string) *Color {
	return color.New(LSColorAttributes[key]...)
}

// NewEXAColor will return `*color.Color` using `EXAColorAttributes[key]`
func NewEXAColor(key string) *Color {
	return color.New(EXAColorAttributes[key]...)
}

// FileLSColor returns color of file (fullpath) according to LS_COLORS
func FileLSColor(fullpath string) *Color {
	fi, err := os.Lstat(fullpath)
	if err != nil {
		return Cerror
	}

	if fi.IsDir() { // os.ModeDir
		return Cdip
	}

	if fi.Mode()&os.ModeSymlink != 0 { // os.ModeSymlink
		// _, err := filepath.EvalSymlinks(fullpath)
		_, err := os.Readlink(fullpath)
		if err != nil {
			return Corp //NewLSColor("or")
		}
		return Clnp
	}

	if fi.Mode()&os.ModeCharDevice != 0 { // os.ModeDevice | os.ModeCharDevice
		return Ccdp
	}

	if fi.Mode()&os.ModeDevice != 0 { //
		return Cbdp
	}

	if fi.Mode()&os.ModeNamedPipe != 0 { //os.ModeNamedPipe
		return Cpip
	}
	if fi.Mode()&os.ModeSocket != 0 { //os.ModeSocket
		return Csop
	}

	// the file is executable by any of its owner, the group and others, use bitmask 0111
	if fi.Mode()&0111 != 0 && !fi.IsDir() {
		return Cexp
	}

	base := filepath.Base(fullpath)
	if att, ok := LSColorAttributes[base]; ok {
		return color.New(att...)
	}
	ext := filepath.Ext(fullpath)
	if att, ok := LSColorAttributes[ext]; ok {
		return color.New(att...)
	}
	file := strings.TrimSuffix(base, ext)
	if att, ok := LSColorAttributes[file]; ok {
		return color.New(att...)
	}
	for re, att := range ReExtLSColors {
		if re.MatchString(base) {
			return color.New(att...)
		}
	}
	return Cfip
}

// FgGray return foreground gray color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of level:
//  level <-> 256 color code
//   0-24 <->  232-255：從黑到白的24階灰度色
func FgGray(level int) *Color {
	return color.New(FgGrayA(level)...)
}

// BgGray return background gray color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of level:
//  level <-> 256 color code
//   0-24 <->  232-255：從黑到白的24階灰度色
func BgGray(level int) *Color {
	return color.New(FgGrayA(level)...)
}
func getGrayCode(level int) (code int) {
	switch i := level; {
	case i < 0:
		code = 0
	case i > 23:
		code = 23
	default:
		code = level
	}
	return code
}

// FgGray return foreground gray color AttributeA
// 	code must be type of int or color.Attribute
// 	range of level:
//  level <-> 256 color code
//   0-24 <->  232-255：從黑到白的24階灰度色
func FgGrayA(level int) AttributeA {
	code := getGrayCode(level)
	return FgColor256A(Grays[code])
}

// BgGray return background gray color AttributeA
// 	code must be type of int or color.Attribute
// 	range of level:
//  level <-> 256 color code
//   0-24 <->  232-255：從黑到白的24階灰度色
func BgGrayA(level int) AttributeA {
	code := getGrayCode(level)
	return BgColor256A(Grays[code])
}

// Color256 return foreground Color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of code:
//      0-  7：標準顏色（同ESC [ 30–37 m）
//      8- 15：高强度颜色（同ESC [ 90–97 m）
//     16-231：6 × 6 × 6 立方（216色）: 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//    232-255：從黑到白的24階灰度色
func FgColor256(code interface{}) *Color {
	return color.New(FgColor256A(code)...)
}

// Color256 return attributes of foreground Color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of code:
//      0-  7：標準顏色（同ESC [ 30–37 m）
//      8- 15：高强度颜色（同ESC [ 90–97 m）
//     16-231：6 × 6 × 6 立方（216色）: 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//    232-255：從黑到白的24階灰度色
func FgColor256A(code interface{}) AttributeA {
	a := getColorAttribute(code)
	as := []Attribute{38, 5}
	as = append(as, a)
	return as
}

func getColorAttribute(code interface{}) (a Attribute) {
	x := indirect(code)
	switch x.(type) {
	case int:
		switch idx := x.(int); {
		case idx < 0:
			a = Attribute(0)
		case idx > 255:
			a = Attribute(255)
		default:
			a = Attribute(idx)
		}
	case Attribute:
		a = x.(Attribute)
	}
	return a
}

// Color256 return background Color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of code:
//      0-  7：標準顏色（同ESC [ 30–37 m）
//      8- 15：高强度颜色（同ESC [ 90–97 m）
//     16-231：6 × 6 × 6 立方（216色）: 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//    232-255：從黑到白的24階灰度色
func BgColor256(code interface{}) *Color {
	return color.New(BgColor256A(code)...)
}

// Color256 return attributes of background Color (use fatih.color)
// 	code must be type of int or color.Attribute
// 	range of code:
//      0-  7：標準顏色（同ESC [ 30–37 m）
//      8- 15：高强度颜色（同ESC [ 90–97 m）
//     16-231：6 × 6 × 6 立方（216色）: 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//    232-255：從黑到白的24階灰度色
func BgColor256A(code interface{}) AttributeA {
	a := getColorAttribute(code)
	as := []Attribute{48, 5}
	as = append(as, a)
	return as
}

type AttributeA []Attribute

func NewAttributeA() AttributeA {
	return make(AttributeA, 0)
}

func (a AttributeA) Add(p ...Attribute) AttributeA {
	a = append(a, p...)
	return a
}
