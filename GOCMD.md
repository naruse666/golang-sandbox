# Goコマンド
## 実行系
### `go run`
実行ファイルは生成されず(一時ファイルで実行後削除)インタプリタのように実行できる。

### `go build`
実行ファイルが生成される。

### `go mod init <module name>`
`go.mod`ファイルが生成される。

### `go mod tidy`
必要なライブラリのインストール、不要ファイルの削除をする。

### `go install <path@version>`
githubなどのリポジトリを参照し、ダウンロードとコンパイルを行い、`$GOPATH/bin`にインストールする。

## フォーマット
### `go fmt`
基本的なフォーマットを行う。

### `staticcheck`
[staticcheck](https://staticcheck.io/)(linter).100%指摘が正しいとは言えないが、向き合う必要がある。  

### `go vet`
構文は合っているが、想定通りの動作をしないものを検知できる。

### `golangci-lint`
[golangci-lint](https://golangci-lint.run/).
複数のツールを実行できる。  
`golangci-lint run`で実行。  
設定ファイルは、`.golangci.(yaml|yml|toml|json)`を選べる。