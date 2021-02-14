package paw

import (
	"os"

	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cast"
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
	EXAColors = map[string][]color.Attribute{
		"ca": LSColors["ca"],
		"cd": LSColors["cd"],
		"di": LSColors["di"],
		"do": LSColors["do"],
		"ex": LSColors["ex"],
		"pi": LSColors["pi"],
		"fi": LSColors["fi"],
		"ln": LSColors["ln"],
		"mh": LSColors["mh"],
		"no": LSColors["no"],
		"or": LSColors["or"],
		"ow": LSColors["ow"],
		"sg": LSColors["sg"],
		"su": LSColors["su"],
		"so": LSColors["so"],
		"st": LSColors["st"],
		"bd": LSColors["bd"],
		"rc": LSColors["rc"],
		// "ur": LSColors["ex"],
		"ur": {38, 5, 230, 1},    // user +r bit
		"uw": {38, 5, 209, 1},    // user +w bit
		"ux": {38, 5, 156, 1, 4}, // user +x bit (files)
		"ue": {38, 5, 156, 1},    // user +x bit (file types)
		"gr": {38, 5, 230, 1},    // group +r bit
		"gw": {38, 5, 209, 1},    // group +w bit
		"gx": {38, 5, 156, 1, 4}, // group +x bit
		"tr": {38, 5, 230, 1},    // others +r bit
		"tw": {38, 5, 209, 1},    // others +w bit
		"tx": {38, 5, 156, 1, 4}, // others +x bit
		"sn": {38, 5, 156, 1},    // size number
		"sb": {38, 5, 156},       // size unit
		"uu": {38, 5, 229, 1},    // user is you + 1 -> bold
		// "un": {38, 5, 214},    // user is not you
		"un": {38, 5, 251},    // user is not you
		"gu": {38, 5, 229, 1}, // group with you in it
		// "gn": {38, 5, 214},    // group without you
		"gn": {38, 5, 251}, // group without you
		"da": {38, 5, 153}, // timestamp + 8 -> concealed
		// "hd": {4, 38, 5, 15}, // head
		"hd":  {38, 5, 251, 4, 48, 5, 236}, // head + 4-> underline
		"-":   {38, 5, 8},                  // Concealed
		".":   {38, 5, 8},                  // Concealed
		" ":   {38, 5, 8},                  // Concealed
		"ga":  {38, 5, 156},                // git new
		"gm":  {38, 5, 117},                // git modified
		"gd":  {38, 5, 209},                // git deleted
		"gv":  {38, 5, 230},                // git renamed
		"gt":  {38, 5, 135},                // git type change
		"dir": {38, 5, 189},                //addition 'dir'
		// "xattr": {38, 5, 249, 4}, //addition 'xattr'+ 4-> underline
		"xattr":    {38, 5, 8, 4, 48, 5, 234},
		"xsymb":    {38, 5, 8, 48, 5, 234},
		"in":       {38, 5, 213},    // inode
		"lk":       {38, 5, 209, 1}, // links
		"bk":       {38, 5, 189},    // blocks
		"prompt":   {38, 5, 251, 48, 5, 236},
		"promptsn": {38, 5, 156, 1, 48, 5, 236},
		"promptsu": {38, 5, 156, 48, 5, 236},
		"trace":    {38, 5, color.FgCyan},
		"debug":    {38, 5, color.FgMagenta},
		"info":     {38, 5, color.FgBlue},
		"warn":     {38, 5, color.FgHiRed},
		"error":    {38, 5, 220, 1, 48, 5, 160},
		"fatal":    {38, 5, 220, 1, 48, 5, 160},
		"panic":    {38, 5, 220, 1, 48, 5, 160},
		"md5":      LSColors["no"], //LSColors[".md5"],
	}
	// LSColors = make(map[string]string) is LS_COLORS code according to
	// extention of file
	LSColors = map[string][]color.Attribute{
		"bd":                   {38, 5, 68},
		"ca":                   {38, 5, 17},
		"cd":                   {38, 5, 113, 1},
		"di":                   {38, 5, 30},
		"do":                   {38, 5, 127},
		"ex":                   {38, 5, 208, 1},
		"pi":                   {38, 5, 126},
		"fi":                   {0},
		"ln":                   {38, 5, 45},
		"mh":                   {38, 5, 222, 1},
		"no":                   {0},
		"or":                   {48, 5, 196, 38, 5, 232, 1},
		"ow":                   {38, 5, 220, 1},
		"sg":                   {48, 5, 3, 38, 5, 0},
		"su":                   {38, 5, 220, 1, 3, 100, 1},
		"so":                   {38, 5, 197},
		"st":                   {38, 5, 86, 48, 5, 234},
		"tw":                   {48, 5, 235, 38, 5, 139, 3},
		"LS_COLORS":            {48, 5, 89, 38, 5, 197, 1, 3, 4, 7},
		"-":                    {38, 5, 8}, // Concealed
		".":                    {38, 5, 8}, // Concealed
		"README":               {38, 5, 220, 1, 4},
		"README.rst":           {38, 5, 220, 1, 4},
		"README.md":            {38, 5, 220, 1, 4},
		"LICENSE":              {38, 5, 220, 1, 4},
		"COPYING":              {38, 5, 220, 1, 4},
		"INSTALL":              {38, 5, 220, 1, 4},
		"COPYRIGHT":            {38, 5, 220, 1, 4},
		"AUTHORS":              {38, 5, 220, 1, 4},
		"HISTORY":              {38, 5, 220, 1, 4},
		"CONTRIBUTORS":         {38, 5, 220, 1, 4},
		"PATENTS":              {38, 5, 220, 1, 4},
		"VERSION":              {38, 5, 220, 1, 4},
		"NOTICE":               {38, 5, 220, 1, 4},
		"CHANGES":              {38, 5, 220, 1, 4},
		".log":                 {38, 5, 190},
		".txt":                 {38, 5, 253},
		".etx":                 {38, 5, 184},
		".info":                {38, 5, 184},
		".markdown":            {38, 5, 184},
		".md":                  {38, 5, 184},
		".mkd":                 {38, 5, 184},
		".nfo":                 {38, 5, 184},
		".pod":                 {38, 5, 184},
		".rst":                 {38, 5, 184},
		".tex":                 {38, 5, 184},
		".textile":             {38, 5, 184},
		".bib":                 {38, 5, 178},
		".json":                {38, 5, 178},
		".jsonl":               {38, 5, 178},
		".jsonnet":             {38, 5, 178},
		".libsonnet":           {38, 5, 142},
		".ndjson":              {38, 5, 178},
		".msg":                 {38, 5, 178},
		".pgn":                 {38, 5, 178},
		".rss":                 {38, 5, 178},
		".xml":                 {38, 5, 178},
		".fxml":                {38, 5, 178},
		".toml":                {38, 5, 178},
		".yaml":                {38, 5, 178},
		".yml":                 {38, 5, 178},
		".RData":               {38, 5, 178},
		".rdata":               {38, 5, 178},
		".xsd":                 {38, 5, 178},
		".dtd":                 {38, 5, 178},
		".sgml":                {38, 5, 178},
		".rng":                 {38, 5, 178},
		".rnc":                 {38, 5, 178},
		".cbr":                 {38, 5, 141},
		".cbz":                 {38, 5, 141},
		".chm":                 {38, 5, 141},
		".djvu":                {38, 5, 141},
		".pdf":                 {38, 5, 141},
		".PDF":                 {38, 5, 141},
		".mobi":                {38, 5, 141},
		".epub":                {38, 5, 141},
		".docm":                {38, 5, 111, 4},
		".doc":                 {38, 5, 111},
		".docx":                {38, 5, 111},
		".odb":                 {38, 5, 111},
		".odt":                 {38, 5, 111},
		".rtf":                 {38, 5, 111},
		".odp":                 {38, 5, 166},
		".pps":                 {38, 5, 166},
		".ppt":                 {38, 5, 166},
		".pptx":                {38, 5, 166},
		".ppts":                {38, 5, 166},
		".pptxm":               {38, 5, 166, 4},
		".pptsm":               {38, 5, 166, 4},
		".csv":                 {38, 5, 78},
		".tsv":                 {38, 5, 78},
		".ods":                 {38, 5, 112},
		".xla":                 {38, 5, 76},
		".xls":                 {38, 5, 112},
		".xlsx":                {38, 5, 112},
		".xlsxm":               {38, 5, 112, 4},
		".xltm":                {38, 5, 73, 4},
		".xltx":                {38, 5, 73},
		".pages":               {38, 5, 111},
		".numbers":             {38, 5, 112},
		".key":                 {38, 5, 166},
		"config":               {1},
		"cfg":                  {1},
		"conf":                 {1},
		"rc":                   {1},
		"authorized_keys":      {1},
		"known_hosts":          {1},
		".ini":                 {1},
		".plist":               {1},
		".viminfo":             {1},
		".pcf":                 {1},
		".psf":                 {1},
		".hidden-color-scheme": {1},
		".hidden-tmTheme":      {1},
		".last-run":            {1},
		".merged-ca-bundle":    {1},
		".sublime-build":       {1},
		".sublime-commands":    {1},
		".sublime-keymap":      {1},
		".sublime-settings":    {1},
		".sublime-snippet":     {1},
		".sublime-project":     {1},
		".sublime-workspace":   {1},
		".tmTheme":             {1},
		".user-ca-bundle":      {1},
		".epf":                 {1},
		".git":                 {38, 5, 197},
		".gitignore":           {38, 5, 240},
		".gitattributes":       {38, 5, 240},
		".gitmodules":          {38, 5, 240},
		".awk":                 {38, 5, 172},
		".bash":                {38, 5, 172},
		".bat":                 {38, 5, 172},
		".BAT":                 {38, 5, 172},
		".sed":                 {38, 5, 172},
		".sh":                  {38, 5, 172},
		".zsh":                 {38, 5, 172},
		".vim":                 {38, 5, 172},
		".kak":                 {38, 5, 172},
		".ahk":                 {38, 5, 41},
		".py":                  {38, 5, 41},
		".ipynb":               {38, 5, 41},
		".rb":                  {38, 5, 41},
		".gemspec":             {38, 5, 41},
		".pl":                  {38, 5, 208},
		".PL":                  {38, 5, 160},
		".t":                   {38, 5, 114},
		".msql":                {38, 5, 222},
		".mysql":               {38, 5, 222},
		".pgsql":               {38, 5, 222},
		".sql":                 {38, 5, 222},
		".tcl":                 {38, 5, 64, 1},
		".r":                   {38, 5, 49},
		".R":                   {38, 5, 49},
		".gs":                  {38, 5, 81},
		".clj":                 {38, 5, 41},
		".cljs":                {38, 5, 41},
		".cljc":                {38, 5, 41},
		".cljw":                {38, 5, 41},
		".scala":               {38, 5, 41},
		".sc":                  {38, 5, 41},
		".dart":                {38, 5, 51},
		".asm":                 {38, 5, 81},
		".cl":                  {38, 5, 81},
		".lisp":                {38, 5, 81},
		".rkt":                 {38, 5, 81},
		".lua":                 {38, 5, 81},
		".moon":                {38, 5, 81},
		".c":                   {38, 5, 81},
		".C":                   {38, 5, 81},
		".h":                   {38, 5, 110},
		".H":                   {38, 5, 110},
		".tcc":                 {38, 5, 110},
		".c++":                 {38, 5, 81},
		".h++":                 {38, 5, 110},
		".hpp":                 {38, 5, 110},
		".hxx":                 {38, 5, 110},
		".ii":                  {38, 5, 110},
		".M":                   {38, 5, 110},
		".m":                   {38, 5, 110},
		".cc":                  {38, 5, 81},
		".cs":                  {38, 5, 81},
		".cp":                  {38, 5, 81},
		".cpp":                 {38, 5, 81},
		".cxx":                 {38, 5, 81},
		".cr":                  {38, 5, 81},
		".go":                  {38, 5, 81},
		".f":                   {38, 5, 81},
		".F":                   {38, 5, 81},
		".for":                 {38, 5, 81},
		".ftn":                 {38, 5, 81},
		".f90":                 {38, 5, 81},
		".F90":                 {38, 5, 81},
		".f95":                 {38, 5, 81},
		".F95":                 {38, 5, 81},
		".f03":                 {38, 5, 81},
		".F03":                 {38, 5, 81},
		".f08":                 {38, 5, 81},
		".F08":                 {38, 5, 81},
		".nim":                 {38, 5, 81},
		".nimble":              {38, 5, 81},
		".s":                   {38, 5, 110},
		".S":                   {38, 5, 110},
		".rs":                  {38, 5, 81},
		".scpt":                {38, 5, 219},
		".swift":               {38, 5, 219},
		".sx":                  {38, 5, 81},
		".vala":                {38, 5, 81},
		".vapi":                {38, 5, 81},
		".hi":                  {38, 5, 110},
		".hs":                  {38, 5, 81},
		".lhs":                 {38, 5, 81},
		".agda":                {38, 5, 81},
		".lagda":               {38, 5, 81},
		".lagda.tex":           {38, 5, 81},
		".lagda.rst":           {38, 5, 81},
		".lagda.md":            {38, 5, 81},
		".agdai":               {38, 5, 110},
		".zig":                 {38, 5, 81},
		".v":                   {38, 5, 81},
		".pyc":                 {38, 5, 240},
		".tf":                  {38, 5, 168},
		".tfstate":             {38, 5, 168},
		".tfvars":              {38, 5, 168},
		".css":                 {38, 5, 125, 1},
		".less":                {38, 5, 125, 1},
		".sass":                {38, 5, 125, 1},
		".scss":                {38, 5, 125, 1},
		".htm":                 {38, 5, 125, 1},
		".html":                {38, 5, 125, 1},
		".jhtm":                {38, 5, 125, 1},
		".mht":                 {38, 5, 125, 1},
		".eml":                 {38, 5, 125, 1},
		".mustache":            {38, 5, 125, 1},
		".coffee":              {38, 5, 074, 1},
		".java":                {38, 5, 074, 1},
		".js":                  {38, 5, 074, 1},
		".mjs":                 {38, 5, 074, 1},
		".jsm":                 {38, 5, 074, 1},
		".jsp":                 {38, 5, 074, 1},
		".php":                 {38, 5, 81},
		".ctp":                 {38, 5, 81},
		".twig":                {38, 5, 81},
		".vb":                  {38, 5, 81},
		".vba":                 {38, 5, 81},
		".vbs":                 {38, 5, 81},
		"Dockerfile":           {38, 5, 155, 4},
		".dockerignore":        {38, 5, 240},
		"Makefile":             {38, 5, 155, 4},
		"MANIFEST":             {38, 5, 243, 4},
		"pm_to_blib":           {38, 5, 240},
		".nix":                 {38, 5, 155},
		".dhall":               {38, 5, 178},
		".rake":                {38, 5, 155},
		".am":                  {38, 5, 242},
		".in":                  {38, 5, 242},
		".hin":                 {38, 5, 242},
		".scan":                {38, 5, 242},
		".m4":                  {38, 5, 242},
		".old":                 {38, 5, 242},
		".out":                 {38, 5, 242},
		".SKIP":                {38, 5, 244},
		".diff":                {48, 5, 197, 38, 5, 232},
		".patch":               {48, 5, 197, 38, 5, 232, 1},
		".bmp":                 {38, 5, 97},
		".dicom":               {38, 5, 97},
		".tiff":                {38, 5, 97},
		".tif":                 {38, 5, 97},
		".TIFF":                {38, 5, 97},
		".cdr":                 {38, 5, 97},
		".flif":                {38, 5, 97},
		".gif":                 {38, 5, 97},
		".icns":                {38, 5, 97},
		".ico":                 {38, 5, 97},
		".jpeg":                {38, 5, 97},
		".JPG":                 {38, 5, 97},
		".jpg":                 {38, 5, 97},
		".nth":                 {38, 5, 97},
		".png":                 {38, 5, 97},
		".psd":                 {38, 5, 97},
		".pxd":                 {38, 5, 97},
		".pxm":                 {38, 5, 97},
		".xpm":                 {38, 5, 97},
		".webp":                {38, 5, 97},
		".ai":                  {38, 5, 99},
		".eps":                 {38, 5, 99},
		".epsf":                {38, 5, 99},
		".drw":                 {38, 5, 99},
		".ps":                  {38, 5, 99},
		".svg":                 {38, 5, 99},
		".avi":                 {38, 5, 114},
		".divx":                {38, 5, 114},
		".IFO":                 {38, 5, 114},
		".m2v":                 {38, 5, 114},
		".m4v":                 {38, 5, 114},
		".mkv":                 {38, 5, 114},
		".MOV":                 {38, 5, 114},
		".mov":                 {38, 5, 114},
		".mp4":                 {38, 5, 114},
		".mpeg":                {38, 5, 114},
		".mpg":                 {38, 5, 114},
		".ogm":                 {38, 5, 114},
		".rmvb":                {38, 5, 114},
		".sample":              {38, 5, 114},
		".wmv":                 {38, 5, 114},
		".3g2":                 {38, 5, 115},
		".3gp":                 {38, 5, 115},
		".gp3":                 {38, 5, 115},
		".webm":                {38, 5, 115},
		".gp4":                 {38, 5, 115},
		".asf":                 {38, 5, 115},
		".flv":                 {38, 5, 115},
		".ts":                  {38, 5, 115},
		".ogv":                 {38, 5, 115},
		".f4v":                 {38, 5, 115},
		".VOB":                 {38, 5, 115, 1},
		".vob":                 {38, 5, 115, 1},
		".ass":                 {38, 5, 117},
		".srt":                 {38, 5, 117},
		".ssa":                 {38, 5, 117},
		".sub":                 {38, 5, 117},
		".sup":                 {38, 5, 117},
		".vtt":                 {38, 5, 117},
		".3ga":                 {38, 5, 137, 1},
		".S3M":                 {38, 5, 137, 1},
		".aac":                 {38, 5, 137, 1},
		".amr":                 {38, 5, 137, 1},
		".au":                  {38, 5, 137, 1},
		".caf":                 {38, 5, 137, 1},
		".dat":                 {38, 5, 137, 1},
		".dts":                 {38, 5, 137, 1},
		".fcm":                 {38, 5, 137, 1},
		".m4a":                 {38, 5, 137, 1},
		".mod":                 {38, 5, 137, 1},
		".mp3":                 {38, 5, 137, 1},
		".mp4a":                {38, 5, 137, 1},
		".oga":                 {38, 5, 137, 1},
		".ogg":                 {38, 5, 137, 1},
		".opus":                {38, 5, 137, 1},
		".s3m":                 {38, 5, 137, 1},
		".sid":                 {38, 5, 137, 1},
		".wma":                 {38, 5, 137, 1},
		".ape":                 {38, 5, 136, 1},
		".aiff":                {38, 5, 136, 1},
		".cda":                 {38, 5, 136, 1},
		".flac":                {38, 5, 136, 1},
		".alac":                {38, 5, 136, 1},
		".mid":                 {38, 5, 136, 1},
		".midi":                {38, 5, 136, 1},
		".pcm":                 {38, 5, 136, 1},
		".wav":                 {38, 5, 136, 1},
		".wv":                  {38, 5, 136, 1},
		".wvc":                 {38, 5, 136, 1},
		".afm":                 {38, 5, 66},
		".fon":                 {38, 5, 66},
		".fnt":                 {38, 5, 66},
		".pfb":                 {38, 5, 66},
		".pfm":                 {38, 5, 66},
		".ttf":                 {38, 5, 66},
		".otf":                 {38, 5, 66},
		".woff":                {38, 5, 66},
		".woff2":               {38, 5, 66},
		".PFA":                 {38, 5, 66},
		".pfa":                 {38, 5, 66},
		".7z":                  {38, 5, 40},
		".a":                   {38, 5, 40},
		".arj":                 {38, 5, 40},
		".bz2":                 {38, 5, 40},
		".cpio":                {38, 5, 40},
		".gz":                  {38, 5, 40},
		".lrz":                 {38, 5, 40},
		".lz":                  {38, 5, 40},
		".lzma":                {38, 5, 40},
		".lzo":                 {38, 5, 40},
		".rar":                 {38, 5, 40},
		".s7z":                 {38, 5, 40},
		".sz":                  {38, 5, 40},
		".tar":                 {38, 5, 40},
		".tgz":                 {38, 5, 40},
		".warc":                {38, 5, 40},
		".WARC":                {38, 5, 40},
		".xz":                  {38, 5, 40},
		".z":                   {38, 5, 40},
		".zip":                 {38, 5, 40},
		".zipx":                {38, 5, 40},
		".zoo":                 {38, 5, 40},
		".zpaq":                {38, 5, 40},
		".zst":                 {38, 5, 40},
		".zstd":                {38, 5, 40},
		".zz":                  {38, 5, 40},
		".apk":                 {38, 5, 215},
		".ipa":                 {38, 5, 215},
		".deb":                 {38, 5, 215},
		".rpm":                 {38, 5, 215},
		".jad":                 {38, 5, 215},
		".jar":                 {38, 5, 215},
		".cab":                 {38, 5, 215},
		".pak":                 {38, 5, 215},
		".pk3":                 {38, 5, 215},
		".vdf":                 {38, 5, 215},
		".vpk":                 {38, 5, 215},
		".bsp":                 {38, 5, 215},
		".dmg":                 {38, 5, 215},
		".r[0-9]{0,2}":         {38, 5, 239},
		".zx[0-9]{0,2}":        {38, 5, 239},
		".z[0-9]{0,2}":         {38, 5, 239},
		".part":                {38, 5, 239},
		".iso":                 {38, 5, 124},
		".bin":                 {38, 5, 124},
		".nrg":                 {38, 5, 124},
		".qcow":                {38, 5, 124},
		".sparseimage":         {38, 5, 124},
		".toast":               {38, 5, 124},
		".vcd":                 {38, 5, 124},
		".vmdk":                {38, 5, 124},
		".accdb":               {38, 5, 60},
		".accde":               {38, 5, 60},
		".accdr":               {38, 5, 60},
		".accdt":               {38, 5, 60},
		".db":                  {38, 5, 60},
		".fmp12":               {38, 5, 60},
		".fp7":                 {38, 5, 60},
		".localstorage":        {38, 5, 60},
		".mdb":                 {38, 5, 60},
		".mde":                 {38, 5, 60},
		".sqlite":              {38, 5, 60},
		".typelib":             {38, 5, 60},
		".nc":                  {38, 5, 60},
		".pacnew":              {38, 5, 33},
		".un~":                 {38, 5, 241},
		".orig":                {38, 5, 241},
		".BUP":                 {38, 5, 241},
		".bak":                 {38, 5, 241},
		".o":                   {38, 5, 241},
		"core":                 {38, 5, 241},
		".mdump":               {38, 5, 241},
		".rlib":                {38, 5, 241},
		".dll":                 {38, 5, 241},
		".swp":                 {38, 5, 244},
		".swo":                 {38, 5, 244},
		".tmp":                 {38, 5, 244},
		".sassc":               {38, 5, 244},
		".pid":                 {38, 5, 248},
		".state":               {38, 5, 248},
		"lockfile":             {38, 5, 248},
		"lock":                 {38, 5, 248},
		".err":                 {38, 5, 160, 1},
		".error":               {38, 5, 160, 1},
		".stderr":              {38, 5, 160, 1},
		".aria2":               {38, 5, 241},
		".dump":                {38, 5, 241},
		".stackdump":           {38, 5, 241},
		".zcompdump":           {38, 5, 241},
		".zwc":                 {38, 5, 241},
		".pcap":                {38, 5, 29},
		".cap":                 {38, 5, 29},
		".dmp":                 {38, 5, 29},
		".DS_Store":            {38, 5, 239},
		".localized":           {38, 5, 239},
		".CFUserTextEncoding":  {38, 5, 239},
		".allow":               {38, 5, 112},
		".deny":                {38, 5, 196},
		".service":             {38, 5, 45},
		"@.service":            {38, 5, 45},
		".socket":              {38, 5, 45},
		".swap":                {38, 5, 45},
		".device":              {38, 5, 45},
		".mount":               {38, 5, 45},
		".automount":           {38, 5, 45},
		".target":              {38, 5, 45},
		".path":                {38, 5, 45},
		".timer":               {38, 5, 45},
		".snapshot":            {38, 5, 45},
		".application":         {38, 5, 116},
		".cue":                 {38, 5, 116},
		".description":         {38, 5, 116},
		".directory":           {38, 5, 116},
		".m3u":                 {38, 5, 116},
		".m3u8":                {38, 5, 116},
		".md5":                 {38, 5, 116},
		".properties":          {38, 5, 116},
		".sfv":                 {38, 5, 116},
		".theme":               {38, 5, 116},
		".torrent":             {38, 5, 116},
		".urlview":             {38, 5, 116},
		".webloc":              {38, 5, 116},
		".lnk":                 {38, 5, 39},
		"CodeResources":        {38, 5, 239},
		"PkgInfo":              {38, 5, 239},
		".nib":                 {38, 5, 57},
		".car":                 {38, 5, 57},
		".dylib":               {38, 5, 241},
		".entitlements":        {1},
		".pbxproj":             {1},
		".strings":             {1},
		".storyboard":          {38, 5, 196},
		".xcconfig":            {1},
		".xcsettings":          {1},
		".xcuserstate":         {1},
		".xcworkspacedata":     {1},
		".xib":                 {38, 5, 208},
		".asc":                 {38, 5, 192, 3},
		".bfe":                 {38, 5, 192, 3},
		".enc":                 {38, 5, 192, 3},
		".gpg":                 {38, 5, 192, 3},
		".signature":           {38, 5, 192, 3},
		".sig":                 {38, 5, 192, 3},
		".p12":                 {38, 5, 192, 3},
		".pem":                 {38, 5, 192, 3},
		".pgp":                 {38, 5, 192, 3},
		".p7s":                 {38, 5, 192, 3},
		"id_dsa":               {38, 5, 192, 3},
		"id_rsa":               {38, 5, 192, 3},
		"id_ecdsa":             {38, 5, 192, 3},
		"id_ed25519":           {38, 5, 192, 3},
		".32x":                 {38, 5, 213},
		".cdi":                 {38, 5, 213},
		".fm2":                 {38, 5, 213},
		".rom":                 {38, 5, 213},
		".sav":                 {38, 5, 213},
		".st":                  {38, 5, 213},
		".a00":                 {38, 5, 213},
		".a52":                 {38, 5, 213},
		".A64":                 {38, 5, 213},
		".a64":                 {38, 5, 213},
		".a78":                 {38, 5, 213},
		".adf":                 {38, 5, 213},
		".atr":                 {38, 5, 213},
		".gb":                  {38, 5, 213},
		".gba":                 {38, 5, 213},
		".gbc":                 {38, 5, 213},
		".gel":                 {38, 5, 213},
		".gg":                  {38, 5, 213},
		".ggl":                 {38, 5, 213},
		".ipk":                 {38, 5, 213},
		".j64":                 {38, 5, 213},
		".nds":                 {38, 5, 213},
		".nes":                 {38, 5, 213},
		".sms":                 {38, 5, 213},
		".8xp":                 {38, 5, 121},
		".8eu":                 {38, 5, 121},
		".82p":                 {38, 5, 121},
		".83p":                 {38, 5, 121},
		".8xe":                 {38, 5, 121},
		".stl":                 {38, 5, 216},
		".dwg":                 {38, 5, 216},
		".ply":                 {38, 5, 216},
		".wrl":                 {38, 5, 216},
		".pot":                 {38, 5, 7},
		".pcb":                 {38, 5, 7},
		".mm":                  {38, 5, 7},
		".gbr":                 {38, 5, 7},
		".scm":                 {38, 5, 7},
		".xcf":                 {38, 5, 7},
		".spl":                 {38, 5, 7},
		".Rproj":               {38, 5, 11},
		".sis":                 {38, 5, 7},
		".1p":                  {38, 5, 7},
		".3p":                  {38, 5, 7},
		".cnc":                 {38, 5, 7},
		".def":                 {38, 5, 7},
		".ex":                  {38, 5, 7},
		".example":             {38, 5, 7},
		".feature":             {38, 5, 7},
		".ger":                 {38, 5, 7},
		".ics":                 {38, 5, 7},
		".map":                 {38, 5, 7},
		".mf":                  {38, 5, 7},
		".mfasl":               {38, 5, 7},
		".mi":                  {38, 5, 7},
		".mtx":                 {38, 5, 7},
		".pc":                  {38, 5, 7},
		".pi":                  {38, 5, 7},
		".plt":                 {38, 5, 7},
		".pm":                  {38, 5, 7},
		".rdf":                 {38, 5, 7},
		".ru":                  {38, 5, 7},
		".sch":                 {38, 5, 7},
		".sty":                 {38, 5, 7},
		".sug":                 {38, 5, 7},
		".tdy":                 {38, 5, 7},
		".tfm":                 {38, 5, 7},
		".tfnt":                {38, 5, 7},
		".tg":                  {38, 5, 7},
		".vcard":               {38, 5, 7},
		".vcf":                 {38, 5, 7},
		".xln":                 {38, 5, 7},
		".iml":                 {38, 5, 166},
	}
	// ReExtLSColors is LS_COLORS code for specific pattern of file extentions
	ReExtLSColors = map[*regexp.Regexp][]color.Attribute{
		regexp.MustCompile(`r[0-9]{0,2}$`):  {38, 5, 239},
		regexp.MustCompile(`zx[0-9]{0,2}$`): {38, 5, 239},
		regexp.MustCompile(`z[0-9]{0,2}$`):  {38, 5, 239},
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
	Csup = NewEXAColor("sn")
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
	// Cprompt is default color use for prompt
	Cpmpt = NewEXAColor("prompt")
	// CpmptSn is default color use for number in prompt
	CpmptSn = NewEXAColor("promptsn")
	// CpmptSu is default color use for unit in prompt
	CpmptSu = NewEXAColor("promptsu")
	Ctrace  = NewEXAColor("trace")
	Cdebug  = NewEXAColor("debug")
	Cinfo   = NewEXAColor("info")
	Cwarn   = NewEXAColor("warn")
	Cerror  = NewEXAColor("error")
	Cfatal  = NewEXAColor("fatal")
	Cpanic  = NewEXAColor("panic")
)

func init() {

}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var reANSI = regexp.MustCompile(ansi)

// StripANSI returns a string without ESC color code
func StripANSI(str string) string {
	return reANSI.ReplaceAllString(str, "")
}

// SetNoColor will set `true` to `NoColor`
func SetNoColor() {
	NoColor = true
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
			LSColors[kv[0]] = getColorAttribute(kv[1])
		}
	}
}

func getColorAttribute(code string) []color.Attribute {
	att := []color.Attribute{}
	for _, a := range strings.Split(code, ";") {
		att = append(att, color.Attribute(cast.ToInt(a)))
	}
	return att
}

// KindLSColorString will colorful string `s` using key `kind`
func KindLSColorString(kind, s string) string {
	att, ok := LSColors[kind]
	if !ok {
		att = LSColors["fi"]
	}
	return color.New(att...).Sprint(s)
}

func KindEXAColorString(kind, s string) string {
	att, ok := EXAColors[kind]
	if !ok {
		att = EXAColors["fi"]
	}
	return color.New(att...).Sprint(s)
}

// func colorstr(att []color.Attribute, s string) string {
// 	cs := color.New(att...)
// 	return cs.Sprint(s)
// }

// // getColorExt will return the color key of extention from `fullpath`
// func getColorExt(fullpath string) (file, ext string) {
// 	fi, err := os.Lstat(fullpath)
// 	if err != nil {
// 		return "", "no"
// 	}
// 	mode := fi.Mode()
// 	sperm := fmt.Sprintf("%v", mode)
// 	switch {
// 	case mode.IsDir(): // d: is a directory 資料夾模式
// 		ext = "di" // di = directory
// 	case mode&os.ModeSymlink != 0: // L: symbolic link 象徵性的關聯
// 		ext = "«link»"
// 		// ext = "ln"
// 		// link, err := filepath.EvalSymlinks(fullpath)
// 		// if err != nil {
// 		// 	ext = "or"
// 		// } else {
// 		// 	_, ext = getColorExt(link)
// 		// }
// 	// } else { // mi = non-existent file pointed to by a symbolic link (visible when you type ls -l)
// 	// 	ext = "mi"
// 	// }
// 	case mode&os.ModeSocket != 0: // S: Unix domain socket Unix 主機 socket
// 		ext = "so" // so = socket file
// 	case mode&os.ModeNamedPipe != 0:
// 		ext = "pi" //pi = fifo file
// 	case mode&os.ModeDevice != 0:
// 		ext = "cd"
// 	// bd = block (buffered) special file
// 	case mode&os.ModeCharDevice != 0:
// 		// cd = character (unbuffered) special file
// 		ext = "cd"
// 	case mode.IsRegular() && !mode.IsDir() && strings.Contains(sperm, "x"):
// 		// ex = file which is executable (ie. has 'x' set in permissions)
// 		ext = "ex"
// 	default: // fi = file
// 		ext = filepath.Ext(fullpath)
// 	}

// 	return filepath.Base(fullpath), ext
// }

// NewLSColor will return `*color.Color` using `LSColors[key]`
func NewLSColor(key string) *color.Color {
	return color.New(LSColors[key]...)
}

// NewEXAColor will return `*color.Color` using `EXAColors[key]`
func NewEXAColor(key string) *color.Color {
	return color.New(EXAColors[key]...)
}
