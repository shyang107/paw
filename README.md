# paw

<!-- TOC -->

- [paw](#paw)
- [Functions](#functions)
    - [io](#io)
        - [`LineCount`](#linecount)
        - [`FileLineCount`](#filelinecount)
        - [`ForEachLine`](#foreachline)
    - [path](#path)
        - [`IsFileExist`](#isfileexist)
        - [`IsDirExists`](#isdirexists)
    - [web](#web)
        - [`GetTitleAndURL`](#gettitleandurl)
        - [`GetTitle`](#gettitle)
        - [`GetURL`](#geturl)
    - [text or string](#text-or-string)
        - [`GetAbbrString`](#getabbrstring)
        - [`CountPlaceHolder`](#countplaceholder)
        - [`HasChineseChar`](#haschinesechar)
    - [Table](#table)

<!-- /TOC -->

# Functions

## io

### `LineCount`

`LineCount` counts the number of `\n` for reader `r`

```language-go
func LineCount(r io.Reader) (int, error)
```

> modify from "github.com/liuzl/goutil"

### `FileLineCount`

`FileLineCount` counts the number of `\n` for file `f`

`f`
: could be `gzip` file or `plain text` file

```language-go
func FileLineCount(f string) (int, error)
```

> modify from "github.com/liuzl/goutil"

### `ForEachLine`

`ForEachLine` higher order function that processes each line of text by callback function.

The last non-empty line of input will be processed even if it has no newline.

`br`
: read from `br` reader

`callback`
: the function used to treatment the every line from `br`

```language-go
func ForEachLine(br *bufio.Reader, callback func(string) error) error
```

> modify from "github.com/liuzl/goutil"

## path

### `IsFileExist`

`IsFileExist` return `true` that `fileName` exist or `false` for not exist

```language-go
func IsFileExist(fileName string) bool
```

### `IsDirExists`

`IsDirExists` return `true` that `dir` is directory or false for not

```language-go
func IsDirExists(dir string) bool
```

## web

### `GetTitleAndURL`

`GetTitleAndURL` get the `title` and `URL` of active tab of the current window of `browser`

`browser`
: - "edge" for "Microsoft Edge" (default)
: - "chrome" for "Google Chrome"

```language-go
func GetTitleAndURL(browser string) (t, u string, err error)
```

### `GetTitle`

`GetTitle` get the `title` of active tab of the current window of `browser`

```language-go
func GetTitle(browser string) (string, error)
```

### `GetURL`

`GetURL` get the `URL` of active tab of the current window of `browser`

```language-go
func GetURL(browser string) (string, error)
```

## text or string

### `GetAbbrString`

`GetAbbrString` return a abbreviation string 'xxx...' of `str` with maximum length `maxlen`.

```language-go
func GetAbbrString(str string, maxlen int) string
```

### `CountPlaceHolder`

`CountPlaceHolder` return `nHan` and `nASCII`

`nHan`
: number of occupied space in terminal for han-character

`nASCII`
: number of occupied space in terminal for ASCII-character

```language-go
func CountPlaceHolder(str string) (nHan int, nASCII int)
```

### `HasChineseChar`

`HasChineseChar` return true for that `str` include chinese character

```language-go
func HasChineseChar(str string) bool
```

## Table

<!-- TODO -->
