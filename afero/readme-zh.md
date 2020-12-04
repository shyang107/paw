---
Author:         Yang, Shuh-Hua
CSS:            x-devonthink-item://DFC0FF5B-5FC1-449E-9CA1-3723B86A8EFB
HTML header:    <script src="x-devonthink-item://95960C49-40AF-4233-8338-BD3E9E55BFA4" ></script>
URL:            https://github.com/chinanf-boy/afero-zh
---

<section class="line-numbers">

|---|:---|
Source|[chinanf-boy/afero-zh: 中文翻译: <afero> Go 的文件系统抽象系统 校对 ✅](https://github.com/chinanf-boy/afero-zh)
Date|29 Mar 2019

<div align="center">
<h1 class="title">spf13/afero<br>Go 的檔案系統抽象系統</h1>
</div>

 [![translate-svg]][translate-list]

<!--[![explain]][source]-->

[explain]: http://llever.com/explain.svg
[source]: https://github.com/chinanf-boy/Source-Explain
[translate-svg]: http://llever.com/translate.svg
[translate-list]: https://github.com/chinanf-boy/chinese-translate-list

[中文](./readme.md) | [english](https://github.com/spf13/afero)

## 校對 ✅

<!-- doc-templite START generated -->
<!-- repo = 'spf13/afero' -->
<!-- commit = 'd40851caa0d747393da1ffb28f7f9d8b4eeffebd' -->
<!-- time = '2018-09-07' -->

翻譯的原文 | 與日期 | 最新更新 | 更多
---|---|---|---
[commit] | ⏰ 2018-09-07 | ![last] | [中文翻譯][translate-list]

[last]: https://img.shields.io/github/last-commit/spf13/afero.svg
[commit]: https://github.com/spf13/afero/tree/d40851caa0d747393da1ffb28f7f9d8b4eeffebd

<!-- doc-templite END generated -->

### 翻譯貢獻

歡迎 👏 勘誤/校對/更新貢獻 😊 [具體貢獻請看](https://github.com/chinanf-boy/chinese-translate-list#貢獻)

## 生活

[If help, **buy** me coffee —— 營養跟不上了，給我來瓶營養快線吧! 💰](https://github.com/chinanf-boy/live-need-money)

---

![afero logo-sm](https://cloud.githubusercontent.com/assets/173412/11490338/d50e16dc-97a5-11e5-8b12-019a300d0fcb.png)

Go 的檔案系統抽象系統

[![Build Status](https://travis-ci.org/spf13/afero.svg)](https://travis-ci.org/spf13/afero) [![Build status](https://ci.appveyor.com/api/projects/status/github/spf13/afero?branch=master&svg=true)](https://ci.appveyor.com/project/spf13/afero) [![GoDoc](https://godoc.org/github.com/spf13/afero?status.svg)](https://godoc.org/github.com/spf13/afero) [![Join the chat at https://gitter.im/spf13/afero](https://badges.gitter.im/Dev%20Chat.svg)](https://gitter.im/spf13/afero?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

---

<a name="toc"></a>
<details>
<summary><kbd class="type-writer">Table of Contents</kbd></summary>

{{TOC}}

</details>

---

# 概觀

Afero 是一個檔案系統框架,提供與任何檔案系統的簡單，統一和通用的互動 API。作為提供介面，型別和方法的抽象層，Afero 具有非常乾淨的介面和簡單的設計，沒有不必要的建構函式或初始化方法。

Afero 也是一個庫，提供一組可互操作的後端檔案系統，輕鬆地使用。同時保留 os 和 ioutil 包的所有功能和優點.

Afero 比單獨使用 os 包提供了顯著的改進，最顯著的是能夠在不依賴磁碟的情況下，建立模擬和測試檔案系統.

它考慮到了，您想使用 OS 包的任何情況，因為它提供了額外的抽象，使得在測試期間可以輕鬆使用記憶體支援的檔案系統。它還增加了對 http 檔案系統的支援，以實現完全的互操作性.

## Afero 特性

- 用於訪問各種檔案系統的唯一的一致 API
- 各種檔案系統型別之間的互操作
- 一組介面，用於鼓勵，和實現後端之間的互操作性
- 跨平臺記憶體支援的檔案原子系統
- 通過組合多個檔案系統，來支援一個組合(聯合)檔案系統
- 修改現有檔案系統的專用後端(只讀， `Regexp` 過濾)
- 一組從 `io`， `ioutil` 和 `hugo` 移植到 `afero` 意識的實用函式

# 使用 Afero

Afero 易於使用，且簡單明瞭.

您可以使用 Afero 的幾種不同方式:

- 單獨使用介面，來定義您自己的檔案系統.
- 包裝為 OS 包.
- 為應用程式的不同部分，定義不同的檔案系統.
- 在測試時，使用 Afero 模擬檔案系統

## 第 1 步:安裝 Afero

首先使用 `go get` 安裝最新版本的庫.

```language-bash
$ go get github.com/spf13/afero
```

接下來在您的應用程式中，包含 Afero.

```language-go
import "github.com/spf13/afero"
```

## 第 2 步:宣告後端

首先定義一個包級變數，並將其設定為指向檔案系統的指標.

```language-go
var AppFs = afero.NewMemMapFs()

or

var AppFs = afero.NewOsFs()
```

重要的是要注意，如果重複呼叫，您將使用一個全新的隔離檔案系統。在 `OsFs` 的情況下，它仍將使用相同的底層檔案系統，但會降低根據需要放入其他檔案系統的能力.

## 第 3 步:像作業系統包一樣使用它

在整個應用程式中，使用您通常會使用的任何功能和方法.

所以，如果我以前的應用有:

```language-go
os.Open('/tmp/foo')
```

我們將其替換為:

```language-go
AppFs.Open('/tmp/foo')
```

`AppFs`是我們上面定義的變數.

## 所有可用功能的列表

檔案系統方法可用:

```language-go
Chmod(name string, mode os.FileMode) : error
Chtimes(name string, atime time.Time, mtime time.Time) : error
Create(name string) : File, error
Mkdir(name string, perm os.FileMode) : error
MkdirAll(path string, perm os.FileMode) : error
Name() : string
Open(name string) : File, error
OpenFile(name string, flag int, perm os.FileMode) : File, error
Remove(name string) : error
RemoveAll(path string) : error
Rename(oldname, newname string) : error
Stat(name string) : os.FileInfo, error
```

檔案介面和方法可用:

```language-go
io.Closer
io.Reader
io.ReaderAt
io.Seeker
io.Writer
io.WriterAt

Name() : string
Readdir(count int) : []os.FileInfo, error
Readdirnames(n int) : []string, error
Stat() : os.FileInfo, error
Sync() : error
Truncate(size int64) : error
WriteString(s string) : ret int, err error
```

在某些應用程式中，定義一個只匯出**檔案系統變數**的新包，就可以從任何地方輕鬆訪問。

## 使用 Afero 的實用功能

Afero 提供了一組函式，使其更易於使用底層檔案系統。這些函式主要來自 `io` & `ioutil`，其中一些是為 Hugo 開發的.

afero 實用程式，支援所有 afero 相容的後端.

實用程式列表包括:

```language-go
DirExists(path string) (bool, error)
Exists(path string) (bool, error)
FileContainsBytes(filename string, subslice []byte) (bool, error)
GetTempDir(subPath string) string
IsDir(path string) (bool, error)
IsEmpty(path string) (bool, error)
ReadDir(dirname string) ([]os.FileInfo, error)
ReadFile(filename string) ([]byte, error)
SafeWriteReader(path string, r io.Reader) (err error)
TempDir(dir, prefix string) (name string, err error)
TempFile(dir, prefix string) (f File, err error)
Walk(root string, walkFn filepath.WalkFunc) error
WriteFile(filename string, data []byte, perm os.FileMode) error
WriteReader(path string, r io.Reader) (err error)
```

有關完整列表，請參閱 [Afero 的 GoDoc](https://godoc.org/github.com/spf13/afero)

這裡是有兩種不同的使用方法。

- 您可以直接呼叫它們，每個函式的第一個引數將是檔案系統，或者

- 您可以宣告一個新`Afero`，一種自定義型別，用於將這些函式繫結到，給定檔案系統的方法.

### 直接呼叫實用程式

```language-go
fs := new(afero.MemMapFs)
f, err := afero.TempFile(fs,"", "ioutil-test")
```

### 通過 Afero 呼叫

```language-go
fs := afero.NewMemMapFs()
afs := &afero.Afero{Fs: fs}
f, err := afs.TempFile("", "ioutil-test")
```

## 使用 Afero 進行測試

使用模擬檔案系統進行測試有很大好處。每次初始化時，它都處於完全空白狀態，無論作業系統如何，都可以輕鬆重現。您可以建立重要內容檔案，檔案訪問速度快，同時還可以避免，刪除臨時檔案，Windows 檔案鎖定等所有煩人的問題。`MemMapFs` 後端非常適合測試.

- 比在磁碟上執行 I/O 操作快得多
- 避免安全和許可權問題
- 更多的控制。`rm -rf /` 將充滿信心
- 測試設定要容易得多
- 無需進行測試清理

實現此目的的一種方法是定義如上所述的變數。在您的應用程式測試期間，這將被設定為 `afero.NewOsFs()`，當然您也可以設為 `afero.NewMemMapFs()`.

每個測試都初始化一個空白的平板記憶體後端並不少見。要做到這一點，在我應用程式程式碼中適當的地方，定義 `appFS = afero.NewOsFs()`。此方法可確保測試與順序無關，並且沒有依賴於早期測試留下的狀態.

然後在我的測試中，我會為每個測試初始化一個新的 `MemMapF`:

```language-go
func TestExist(t *testing.T) {
    appFS := afero.NewMemMapFs()
    // 建立 test 檔案 和 目錄
    appFS.MkdirAll("src/a", 0755)
    afero.WriteFile(appFS, "src/a/b", []byte("file b"), 0644)
    afero.WriteFile(appFS, "src/c", []byte("file c"), 0644)
    name := "src/c"
    _, err := appFS.Stat(name)
    if os.IsNotExist(err) {
        t.Errorf("file \"%s\" does not exist.\n", name)
    }
}
```

# 可用的後端

## 原生作業系統

### OsFs

第一個是圍繞原生 OS 呼叫的包裝器。將它變得非常容易使用,因為所有呼叫都與現有的 OS 呼叫相同。它還使您的程式碼在作業系統，與根據需要使用模擬檔案系統時，變得輕鬆，甚至 *無聊\^\_\^*。

```language-go
appfs := afero.NewOsFs()
appfs.MkdirAll("src/a", 0755))
```

## 記憶體支援儲存

### MemMapFs

Afero 還提供完全原子記憶體支援的檔案系統，非常適合用於模擬，並在不需要保持時，加速不必要的磁碟。它是完全併發的，可以安全地在 go 協程中使用.

```language-go
mm := afero.NewMemMapFs()
mm.MkdirAll("src/a", 0755))
```

#### InMemoryFile

作為 MemMapFs 的一部分，Afero 還提供原子的，完全併發的記憶體支援檔案實現.這可以輕鬆地在其他記憶體支援的檔案系統中使用. 計劃是使用 `InMemoryFile` 新增基數樹記憶體儲存檔案系統.

## 網路介面

### SftpFs

Afero 對安全檔案傳輸協議 (`sftp`) 有實驗性的支援。可用於加密通道上執行檔案操作.

## 過濾後端

### BasePathFs

`BasePathF` 將所有操作限制在 `Fs` 內的給定路徑。在呼叫這個源 `Fs` 之前， `Fs` 操作的給定檔名，會以基本路徑為字首.

```language-go
bp := afero.NewBasePathFs(afero.NewOsFs(), "/base/path")
```

### ReadOnlyFs

源 `Fs` 周圍的薄包裝器，提供只讀的.

```language-go
fs := afero.NewReadOnlyFs(afero.NewOsFs())
_, err := fs.Create("/file.txt")
// err = syscall.EPERM
```

# RegexpFs

對檔名進行過濾後的，任何與傳遞的正規表示式不匹配的檔案，都將被視為不存在。將不會建立與提供的正規表示式不匹配的檔案。目錄不過濾.

```language-go
fs := afero.NewRegexpFs(afero.NewMemMapFs(), regexp.MustCompile(`\.txt$`))
_, err := fs.Create("/file.html")
// err = syscall.ENOENT
```

## HttpFs

Afero 提供了一個 `http` 相容的後端,可以包裝任何現有的後端.

`Http` 包需要稍微特定的 `Open` 版本,它返回一個 `http.File` 型別.

Afero 提供滿足此要求的 `httpFs` 檔案系統。任何 Afero FileSystem 都可以用作 `httpFs`。

```language-go
httpFs := afero.NewHttpFs(<ExistingFS>)
fileserver := http.FileServer(httpFs.Dir(<PATH>)))
http.Handle("/", fileserver)
```

## 複合後端

Afero 提供合成兩個檔案系統 (或更多)，作為單個檔案系統的能力.

### CacheOnReadFs

`CacheOnReadFs` 將懶洋洋地將任何訪問過的檔案，從 `基礎層-base` 複製到 `覆蓋層-overlay` 中。後續讀取將直接從覆蓋層中提取，允許快取持續時間內，請求在覆蓋層中建立的快取。

如果基本檔案系統是可寫的,則對檔案的任何更改，將首先對基礎層進行,然後對覆蓋層進行。而開啟檔案的 Write 呼叫控制,如 `Write()` 或 `Truncate()` 則先到覆蓋層.

要僅將檔案寫入覆蓋層,可以直接使用覆蓋層 `Fs` (而不是通過聯合  ).

在給定 `time.Duration` 快取持續時間內，對該層中的檔案進行快取，快取持續時間為 `0` ，意味著"永遠"，意味著檔案將不會從基礎層重新請求.

只讀的基礎層會讓覆蓋層也是隻讀的，但是當檔案在快取層中不存在 (或過時) 時，仍然將檔案從基礎層複製到覆蓋層.

```language-go
base := afero.NewOsFs()
layer := afero.NewMemMapFs()
ufs := afero.NewCacheOnReadFs(base, layer, 100 * time.Second)
```

### CopyOnWriteFs()

`CopyOnWriteFs` 是一個只讀的基本檔案系統，頂部有一個可寫的層.

`Read` 操作首先檢視覆蓋層，如果沒有找到，將從基礎層提供檔案服務.

只能在覆蓋層中對檔案系統進行更改.

任何僅在基礎中找到的檔案的修改，都會在修改 (包括開啟可寫的檔案) 之前，將檔案複製到覆蓋層.

目前不允許刪除和重新命名，僅存在於基礎層中的檔案。如果檔案在基礎層和覆蓋層中存在，則僅能刪除/重新命名覆蓋層.

```language-go
    base := afero.NewOsFs()
    roBase := afero.NewReadOnlyFs(base)
    ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())

    fh, _ = ufs.Create("/home/test/file2.txt")
    fh.WriteString("This is a test")
    fh.Close()
```

在此示例中，所有寫入操作僅發生在記憶體(MemMapFs)中，基本檔案系統(OsFs)保持不變.

## 期望/可能的後端

以下是我們希望有人可能實現的後端，簡短列表:

- SSH
- ZIP
- TAR
- S3

# 關於該專案

## 這個名字是什麼

Afero 來自拉丁美洲的 Ad-Facere.

**"Ad"**是一個字首,意思是"to".

**"Facere"**是"make 或 do "的根單詞"faciō"的一種形式.

afero 的字面含義是"to make"或"to do",這對於允許製作檔案和目錄,並使用它們進行操作的庫來說非常合適.

與 Afero 具有相同根源的英語單詞是"affair"。Affair 擁有相同的概念,但作為名詞，它意味著"製造或完成的東西"或"特定型別的物體".

與我的其他一些庫 (`hugo`， `cobra`， `viper`) 不同，谷歌一下也不錯.

## 發行說明

- **0.10.0** 2015.12.10
    - 與 Windows 完全相容
    - 介紹 afero 實用函式
    - 測試套件重寫，為跨平臺工作
    - 規範化 MemMapFs 的路徑
    - 將 Sync 新增到檔案介面
    - **破而後立** Walk 和 ReadDir 已更改引數順序
    - 將 MemMapFs 使用的型別移動到子包中
    - 一般錯誤修正和改進
- **0.9.0** 2015.11.05
    - 新的 Walk 函式類似於 filepath.Walk
    - MemMapFs.OpenFile 處理 O_CREATE,O_APPEND,O_TRUNC
    - MemMapFs.Remove 現在真的刪除了該檔案
    - InMemoryFile.Readdir 和 Readdirnames 正常工作
    - InMemoryFile 函式將其鎖定以進行併發訪問
    - 測試套件改進
- **0.8.0** 2014.10.28
    - 第一個公開版
    - 介面感覺準備好，為人們構建使用
    - 介面滿足所有已知用途
    - MemMapFs 通過了大部分 OS 測試套件
    - OsFs 通過了大部分 OS 測試套件

## 貢獻

1. 叉-Fork 吧
2. 建立您的功能分支(`git checkout -b my-new-feature`)
3. 提交你的更改(`git commit -am 'Add some feature'`)
4. 推到分支(`git push origin my-new-feature`)
5. 建立新的 Pull 請求

## 貢獻者

名字沒有特別的順序:

- [spf13](https://github.com/spf13)
- [jaqx0r](https://github.com/jaqx0r)
- [mbertschler](https://github.com/mbertschler)
- [XOR 門](https://github.com/xor-gate)

## 執照

Afero 是在 Apache 2.0 許可下發布的。看到[LICENSE.txt](https://github.com/spf13/afero/blob/master/LICENSE.txt)
