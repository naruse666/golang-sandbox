# golang-sandbox
[golang](https://golang.org/)のキャッチアップリポジトリ。

# 他言語との違い
- Goはあまり数値計算に使われない。(したいなら他の言語を使った方がいい)が、複素数がサポートされている。)
- 継承やオーバーロードなどがない!
- ポインタはある
- 明示的な型変換が絶対
- :=は関数内でのみ

### ゼロ値
宣言されているが値が未割り当ての場合に明示的なゼロ値がある。  
以下はそれぞれの型のゼロ値.
```
bool                     = false
int8~int64, uint8~uint64 = 0
float(32|64)             = 0
string, rune             = "" // 空文字
slice                    = nil
map                      = nil
```

### リテラル
#### 整数リテラル
`1_2_34`の様に数値の間に`_`を書けるが,やるなら`0b_1001_1110`の様にする。
#### 浮動小数点数リテラル
`6.03e23`の様な指数表記もでき、上記同様`_`を書ける。
#### runeリテラル
POSIXのパーミッション値とかで、`0o777`で`rwxrwxrwx`とか。これ以外はあまり使われないらしい。

# ドキュメントの内容をさらっと
## Formatting
`gofmt`コマンドでフォーマット可能。

## Commentary
`/* */ block comments`と`// line comments`が使用可能。

## Names
### Package names
```import "bytes"```でインポート。

### Getters
`GetOwner`ではないので注意。
```
owner := obj.Owner()
if owner != user {
    obj.SetOwner(user)
}
```

### Interface names
メソッド名 + erとして命名する。`Read`メソッドなら`Reader`.

### MixedCaps
Goでは`MixedCaps or mixedCaps`のようにアンダースコアは使わない。

## Semicolons
プログラマが書く必要はない。(コンパイラがやる)  
`break continue fallthrough return ++ -- ) }`の後につけられる。

## Control structures
### If
初期化処理もできる。
```
if err := file.Chmod(0664); err != nil {
    log.Print(err)
    return err
}
```

### For
基本的な構文。
```
sum := 0
for i := 0; i < 10; i++ {
    sum += i
}
```

配列やmapを使用する場合、rangeを使用できる。
```
for key, value := range oldMap {
    newMap[key] = value
}
```

最初のフィールド(key, index等)のみが必要な場合は、セカンドフィールドを削除して表記する。
```
for key := range m {
    if key.expired() {
        delete(m, key)
    }
}
```

二番目のフィールドのみ必要なら、アンダースコアを使用する。
```
sum := 0
for _, value := range array {
    sum += value
}
```

Goは`++, --`がステートメントのため以下の様な表現になる。
```
// Reverse a
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
    a[i], a[j] = a[j], a[i]
}
```

## Functions
複数の値を返せる。
### Named result parameters
入力パラメータと戻り値を指定する。
```func nextInt(b []byte, pos int) (value, nextPos int) {```

### Defer
実行の延期的な表現。以下の例ではCloseがOpenと近くにあることも利点の一つ。
```
// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()  // f.Close will run when we're finished.

    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...) // append is discussed later.
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err  // f will be closed if we return here.
        }
    }
    return string(result), nil // f will be closed if we return here.
}
```

以下の例では`4 3 2 1 0`と出力される。
```
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i)
}
```

ちょっと独特。
```
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func un(s string) {
    fmt.Println("leaving:", s)
}

func a() {
    defer un(trace("a"))
    fmt.Println("in a")
}

func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}

func main() {
    b()
}
```
prints
```
entering: b
in b
entering: a
in a
leaving: a
leaving: b
```

## Data
### Allocation with new
new, makeという組み込み関数がある。  
newはメモリを割り当てるが初期化ではなく0にする。new(T)はTのタイプに0を割り当て、アドレスである*Tを返す。

### Constructors and composite literals
以下の様なコンストラクタを。
```
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := new(File)
    f.fd = fd
    f.name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}
```
簡潔にできる。
```
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := File{fd, name, nil, 0}
    return &f
}
```

`field:name`としてラベルをつける。
```
return &File{fd: fd, name: name}
```

### Allocation with make
`make(T, args)`は`slices, maps, channels`を初期化します。   
ポインタは使用できない。newを使用するか変数に明示的に指定する。  
sliceの例：
```
make([]int, 10, 100)
```

## Arrays
慣例的ではない。基本はsliceを利用する。
```
func Sum(a *[3]float64) (sum float64) {
    for _, v := range *a {
        sum += v
    }
    return
}

array := [...]float64{7.0, 8.5, 9.1}
x := Sum(&array)  // Note the explicit address-of operator
```

## Slices
例えばRead関数ではsliceを引数に持ちます。
```
func (f *File) Read(buf []byte) (n int, err error)
```
次のように最初の32バイトをスライスします。
```
n, err := f.Read(buf[0:32])
```

capを使用してsliceの最大容量取得します。
```
func Append(slice, data []byte) []byte {
    l := len(slice)
    if l + len(data) > cap(slice) {  // reallocate
        // Allocate double what's needed, for future growth.
        newSlice := make([]byte, (l+len(data))*2)
        // The copy function is predeclared and works for any slice type.
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:l+len(data)]
    copy(slice[l:], data)
    return slice
}
```

### Two-dimensional slices
```
type Transform [3][3]float64  // A 3x3 array, really an array of arrays.
type LinesOfText [][]byte     // A slice of byte slices.
```

## Maps
key, valueの型指定する。
```
var timeZone = map[string]int{
    "UTC":  0*60*60,
    "EST": -5*60*60,
    "CST": -6*60*60,
    "MST": -7*60*60,
    "PST": -8*60*60,
}
```

keyがアクセサとなる。
```
offset := timeZone["EST"]
```

`comma ok`という慣例がある。
```
var seconds int
var ok bool
seconds, ok = timeZone[tz]
```

この様な例となる。
```
func offset(tz string) int {
    if seconds, ok := timeZone[tz]; ok {
        return seconds
    }
    log.Println("unknown time zone:", tz)
    return 0
}
```

実際の値を気にせず存在のみの確認は`blank identifier`を使用します。
```
_, present := timeZone[tz]
```

## Printing
```
fmt.Printf("Hello %d\n", 23)
fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
fmt.Println("Hello", 23)
fmt.Println(fmt.Sprint("Hello ", 23))
```

decimalとかはこんな感じで書ける。(`%v`はvalue)
```
fmt.Printf("%v\n", timeZone)  // or just fmt.Println(timeZone)
```

構造体を出力する時、`%+v`は名前付き、`%#v`はGo構文で出力する。
```
type T struct {
    a int
    b float64
    c string
}
t := &T{ 7, -2.35, "abc\tdef" }
fmt.Printf("%v\n", t)
fmt.Printf("%+v\n", t)
fmt.Printf("%#v\n", t)
fmt.Printf("%#v\n", timeZone)
```

print

```
&{7 -2.35 abc   def}
&{a:7 b:-2.35 c:abc     def}
&main.T{a:7, b:-2.35, c:"abc\tdef"}
map[string]int{"CST":-21600, "EST":-18000, "MST":-25200, "PST":-28800, "UTC":0}
```

`%T`はタイプを出力。
```fmt.Printf("%T\n", timeZone) // print: map[string]int```

## Append
定義は以下の様になっており、その下のような感じになる。
```
func append(slice []T, elements ...T) []T
```

```
x := []int{1,2,3}
x = append(x, 4, 5, 6)
fmt.Println(x)
// prints [1 2 3 4 5 6]
```

この様にしても上と同様の結果が得られる。
```
x := []int{1,2,3}
y := []int{4,5,6}
x = append(x, y...)
fmt.Println(x)
```

## Initialization
### Constants
`iota`列挙子を使用して列挙定数を生成する。`=`なことに注意。
```
type ByteSize float64

const (
    _           = iota // ignore first value by assigning to blank identifier
    KB ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```

### Variables
constと同じ感じ。
```
var (
    home   = os.Getenv("HOME")
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)
```

### The init function
init関数はpackageのすべての初期化処理が終了後に実行されます。  
init関数の用途はプログラムの検証等です。
```
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    // gopath may be overridden by --gopath flag on command line.
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```

## Methods
### Pointers vs. Values
メソッドのレシーバに名前付きの型を指定する。
```
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
    // Body exactly the same as the Append function defined above.
}
```

上記はメソッドが更新されたsliceを返す必要があるが、レシーバをポインタで受け取り解消します。
```
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // Body as above, without the return.
    *p = slice
}
```
値メソッドは値とポインタで呼び出せるが、ポインタメソッドはポインタでのみ呼び出し可能。

## The blank identifier
Unixの`/dev/null`に似ている。
### The blank identifier in multiple assignment
エラーだけに着目したい時などに使う。この逆はやるべきでない。
```
if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Printf("%s does not exist\n", path)
}
```

### Unused imports and variables
未使用インポートや変数の警告を消すためにblank identifierを使用する。(なるべく不要なのは入れない。コンパイルが遅くなったりする。)
 ```
 package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

var _ = fmt.Printf // For debugging; delete when done.
var _ io.Reader    // For debugging; delete when done.

func main() {
    fd, err := os.Open("test.go")
    if err != nil {
        log.Fatal(err)
    }
    // TODO: use fd.
    _ = fd
}
```

## Embedding
interfaceの結合的な。
```
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// ReadWriter is the interface that combines the Reader and Writer interfaces.
type ReadWriter interface {
    Reader
    Writer
}
```

## Concurrency
### Share by communicating
通信によってメモリを共有する。スローガン的なの.(?)
```
Do not communicate by sharing memory; instead, share memory by communicating.
```

### Goroutines
関数の完了を通知する方法がないため`channels`を使う方のがいいっぽい。
```go list.Sort()  // run list.Sort concurrently; don't wait for it.```
この様な使い方。
```
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println(message)
    }()  // Note the parentheses - must call the function.
}
```

### Channels
mapの様に`make`を使用して割り当て可能。バッファのデフォルトは0.
```
ci := make(chan int)            // unbuffered channel of integers
cj := make(chan int, 0)         // unbuffered channel of integers
cs := make(chan *os.File, 100)  // buffered channel of pointers to Files
```

一つ上のセクションでchannelsの例を示す。
```
c := make(chan int)  // Allocate a channel.
// Start the sort in a goroutine; when it completes, signal on the channel.
go func() {
    list.Sort()
    c <- 1  // Send a signal; value does not matter.
}()
doSomethingForAWhile()
<-c   // Wait for sort to finish; discard sent value.
```

goroutinesの引数に`req`を設定するなどして、forループのreqをユニークなものとして扱う。
```
func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func(req *Request) {
            process(req)
            <-sem
        }(req)
    }
}
```

一般的なサーバを作成する際、channelを使い終了を待機する例。
```
.func handle(queue chan *Request) {
    for r := range queue {
        process(r)
    }
}

func Serve(clientRequests chan *Request, quit chan bool) {
    // Start handlers
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)
    }
    <-quit  // Wait to be told to exit.
}
```

## Errors
`error`の型がある。
```
// PathError records an error and the operation and
// file path that caused it.
type PathError struct {
    Op string    // "open", "unlink", etc.
    Path string  // The associated file.
    Err error    // Returned by the system call.
}

func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

### Panic
```
// A toy implementation of cube root using Newton's method.
func CubeRoot(x float64) float64 {
    z := x/3   // Arbitrary initial value
    for i := 0; i < 1e6; i++ {
        prevz := z
        z -= (z*z*z-x) / (3*z*z)
        if veryClose(z, prevz) {
            return z
        }
    }
    // A million iterations has not converged; something is wrong.
    panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}
```